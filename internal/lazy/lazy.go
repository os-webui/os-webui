package lazy

import (
	"context"
	"sync"
	"sync/atomic"
)

// Option and funcOption implementations remain exactly as your optimized generic versions.
type Option[T any] interface {
	apply(T)
}

// ============================================================================

type valueOptions[T any] struct {
	ctx   context.Context
	done  func(val T, err error)
	count bool // If true, cancels the background task when all interested waiters leave due to timeout/cancellation.
}

func WithContext[T any](ctx context.Context) Option[*valueOptions[T]] {
	return newFuncOption(func(o *valueOptions[T]) {
		if ctx == nil {
			o.ctx = context.Background()
		} else {
			o.ctx = ctx
		}
	})
}

func WithDone[T any](done func(val T, err error)) Option[*valueOptions[T]] {
	return newFuncOption(func(o *valueOptions[T]) {
		o.done = done
	})
}

// WithAutoCancel activates the dynamic telemetry tracker.
// If all calling contexts timeout or abort, the active execution worker is terminated early.
func WithAutoCancel[T any](enabled bool) Option[*valueOptions[T]] {
	return newFuncOption(func(o *valueOptions[T]) {
		o.count = enabled
	})
}

// ============================================================================

type Value[T any] struct {
	opts valueOptions[T]
	rw   sync.RWMutex
	once *onceResult[T]
}

type onceResult[T any] struct {
	done    chan struct{}
	value   T
	err     error
	waiters atomic.Int64 // Lock-free reference telemetry tracking active requests
	cancel  context.CancelFunc
}

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

func (o *Value[T]) Get(ctx context.Context, f func(context.Context) (T, error)) (T, bool, error) {
	// 1. Pre-flight check
	select {
	case <-ctx.Done():
		var zero T
		return zero, true, ctx.Err()
	case <-o.opts.ctx.Done():
		var zero T
		return zero, true, o.opts.ctx.Err()
	default:
	}

	// 2. Fast-path (RLock)
	o.rw.RLock()
	res := o.once
	if res != nil && o.opts.count {
		res.waiters.Add(1) // Register telemetry under read-lock safety
	}
	o.rw.RUnlock()

	if res != nil {
		return o.await(ctx, res)
	}

	// 3. Slow-path (Lock)
	o.rw.Lock()
	if o.once != nil {
		res = o.once
		if o.opts.count {
			res.waiters.Add(1)
		}
		o.rw.Unlock()
		return o.await(ctx, res)
	}

	// 4. Winner allocation
	bgCtx, bgCancel := context.WithCancel(o.opts.ctx)
	res = &onceResult[T]{
		done:   make(chan struct{}),
		cancel: bgCancel,
	}
	if o.opts.count {
		res.waiters.Store(1) // The winner registers as the primary waiter
	}
	o.once = res
	o.rw.Unlock()

	// 5. Asynchronous Bootstrap Execution
	go func() {
		res.value, res.err = f(bgCtx)

		if res.err != nil {
			o.rw.Lock()
			if o.once == res {
				o.once = nil
			}
			o.rw.Unlock()
		}

		if o.opts.done != nil {
			o.opts.done(res.value, res.err)
		}

		res.cancel() // Housekeeping: clean up context allocations
		close(res.done)
	}()

	return o.await(ctx, res)
}

// await handles localized suspension and manages the dynamic reference counting cleanup block.
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
	if o.opts.count {
		// Atomic decrement. If the counter hits 0, it means EVERY SINGLE WAITER
		// has abandoned the quest due to timeouts. The last one out turns off the lights.
		if res.waiters.Add(-1) == 0 {
			res.cancel() // Interrupt the background task f immediately!
		}
	}
}
