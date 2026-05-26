package lazy

import (
	"context"
	"sync"
	"sync/atomic"
)

// valueOptions bundles the internal generic telemetry configurations.
type valueOptions[T any] struct {
	ctx        context.Context
	done       func(val T, err error)
	autoCancel bool // When true, enables reference-counting to auto-cancel the background task if all waiters leave.
}

// WithContext binds a long-lived supervisor context to the initialization life cycle.
func WithContext[T any](ctx context.Context) Option[*valueOptions[T]] {
	return newFuncOption(func(o *valueOptions[T]) {
		if ctx == nil {
			o.ctx = context.Background()
		} else {
			o.ctx = ctx
		}
	})
}

// WithDone registers an atomic post-initialization callback function.
func WithDone[T any](done func(val T, err error)) Option[*valueOptions[T]] {
	return newFuncOption(func(o *valueOptions[T]) {
		o.done = done
	})
}

// WithAutoCancel activates the dynamic reference telemetry tracker.
func WithAutoCancel[T any](enabled bool) Option[*valueOptions[T]] {
	return newFuncOption(func(o *valueOptions[T]) {
		o.autoCancel = enabled
	})
}

// ============================================================================

// Value guards the lazy initialization of a heavy resource of type T.
type Value[T any] struct {
	opts valueOptions[T]
	rw   sync.RWMutex
	once *onceResult[T]
}

// onceResult acts as an isolated generation state container for a specific initialization attempt.
type onceResult[T any] struct {
	done     chan struct{}
	value    T
	err      error
	waiters  atomic.Int64 // Lock-free tracker counting callers currently waiting for this generation.
	cancel   context.CancelFunc
	canceled atomic.Int64 // 1 indicates this generation has been orphaned/aborted and scheduled for eviction.
}

// New creates a ready-to-use type-safe Lazy Value instance using generic functional options.
func New[T any](opt ...Option[*valueOptions[T]]) *Value[T] {
	val := &Value[T]{
		opts: valueOptions[T]{
			ctx: context.Background(),
		},
	}
	for _, o := range opt {
		o.apply(&val.opts)
	}
	return val
}

// Get ensures that f is called exactly once as long as it succeeds.
func (o *Value[T]) Get(ctx context.Context, f func(context.Context) (T, error)) (T, bool, error) {
	// 1. Pre-flight Check
	select {
	case <-ctx.Done():
		var zero T
		return zero, true, ctx.Err()
	case <-o.opts.ctx.Done():
		var zero T
		return zero, true, o.opts.ctx.Err()
	default:
	}

	// 2. Fast-Path (RLock)
	o.rw.RLock()
	res := o.once

	// A generation is usable if it exists, and either autoCancel is disabled OR it hasn't been flipped to canceled.
	if res != nil && (!o.opts.autoCancel || res.canceled.Load() == 0) {
		if o.opts.autoCancel {
			res.waiters.Add(1)
		}
		o.rw.RUnlock()
		return o.await(ctx, res)
	}
	o.rw.RUnlock()

	// 3. Slow-Path (Lock)
	o.rw.Lock()

	if o.once != nil {
		currentRes := o.once
		if !o.opts.autoCancel || currentRes.canceled.Load() == 0 {
			if o.opts.autoCancel {
				currentRes.waiters.Add(1)
			}
			o.rw.Unlock()
			return o.await(ctx, currentRes)
		}
	}

	// 4. Winner Allocation (Fresh New Generation Setup)
	bgCtx, bgCancel := context.WithCancel(o.opts.ctx)
	newRes := &onceResult[T]{
		done:   make(chan struct{}),
		cancel: bgCancel,
	}
	if o.opts.autoCancel {
		newRes.waiters.Store(1)
	}
	o.once = newRes
	o.rw.Unlock()

	// 5. Asynchronous Task Bootstrap
	go func() {
		// FIX: Must reference 'newRes' explicitly inside the closure fence.
		// Using 'res' from the outer scope points to a stale or nil pointer generation.
		newRes.value, newRes.err = f(bgCtx)

		if newRes.err != nil {
			newRes.canceled.Store(1)
			o.rw.Lock()
			if o.once == newRes {
				o.once = nil
			}
			o.rw.Unlock()
		}

		if o.opts.done != nil {
			o.opts.done(newRes.value, newRes.err)
		}

		newRes.cancel()
		close(newRes.done)
	}()

	return o.await(ctx, newRes)
}

func (o *Value[T]) await(ctx context.Context, res *onceResult[T]) (T, bool, error) {
	select {
	case <-ctx.Done():
		o.handleLeave(res)
		var zero T
		return zero, true, ctx.Err()
	case <-o.opts.ctx.Done():
		o.handleLeave(res)
		var zero T
		return zero, true, o.opts.ctx.Err()
	case <-res.done:
		if res.err != nil {
			return res.value, true, res.err
		}
		return res.value, false, nil
	}
}

func (o *Value[T]) handleLeave(res *onceResult[T]) {
	if o.opts.autoCancel {
		if res.waiters.Add(-1) == 0 {
			res.canceled.Store(1)
			res.cancel()

			o.rw.Lock()
			if o.once == res {
				o.once = nil
			}
			o.rw.Unlock()
		}
	}
}
