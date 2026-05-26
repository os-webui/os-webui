package lazy

import (
	"context"
	"sync"
)

// Option defines a highly reusable generic interface for the Functional Options pattern.
// By passing the configuration struct pointer as the generic type argument,
// this single interface can be reused across different components without creating redundant type-specific options.
type Option[T any] interface {
	apply(T)
}

// ============================================================================

type valueOptions[T any] struct {
	ctx  context.Context
	done func(val T, err error)
}

// WithContext binds a long-lived context (e.g., application or plugin lifecycle) to the value.
// This ensures that even if the caller's transient context is canceled, the background initialization
// task can safely continue executing to completion using this long-lived context.
func WithContext[T any](ctx context.Context) Option[*valueOptions[T]] {
	return newFuncOption(func(o *valueOptions[T]) {
		if ctx == nil {
			o.ctx = context.Background()
		} else {
			o.ctx = ctx
		}
	})
}

// WithDone registers a callback function that executes immediately after the initialization function f completes.
// It is fully type-safe and can be used for telemetry, setting up connection pool metrics, or alerting on failures.
func WithDone[T any](done func(val T, err error)) Option[*valueOptions[T]] {
	return newFuncOption(func(o *valueOptions[T]) {
		o.done = done
	})
}

// ============================================================================

// Value guards the lazy initialization of a resource of type T.
// It resolves 4 major concurrency challenges: high-concurrency protection, lazy-loading,
// instant response to client cancellation, and automatic retry eviction upon failure.
type Value[T any] struct {
	opts valueOptions[T]
	rw   sync.RWMutex
	once *onceResult[T]
}

type onceResult[T any] struct {
	done  chan struct{}
	value T
	err   error
}

// New creates a new type-safe Lazy Value instance using the optimized generic functional options.
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
// It returns (value, shared, error).
// CRITICAL SEMANTIC: shared is false ONLY when the specific caller won the race
// and successfully completed the initialization. For all concurrent waiters, aborted requests,
// or failed generations, shared returns true to protect upper layers (like LRU) from duplicate writes.
func (o *Value[T]) Get(ctx context.Context, f func(context.Context) (T, error)) (T, bool, error) {
	// 1. Top-of-function check: Check context status immediately before acquiring any locks.
	select {
	case <-ctx.Done():
		var zero T
		return zero, true, ctx.Err()
	case <-o.opts.ctx.Done():
		var zero T
		return zero, true, o.opts.ctx.Err()
	default:
	}

	// 2. Fast-path check: Use RLock to check if initialization is already underway or completed.
	o.rw.RLock()
	res := o.once
	o.rw.RUnlock()

	if res != nil {
		// Initialization is active or done. Park concurrent callers on the channel level.
		// Multiple goroutines waiting here will not contend for the mutex lock.
		select {
		case <-ctx.Done():
			var zero T
			return zero, true, ctx.Err()
		case <-o.opts.ctx.Done():
			var zero T
			return zero, true, o.opts.ctx.Err()
		case <-res.done:
			return res.value, true, res.err
		}
	}

	// 3. Slow-path: Upgrade to Write Lock to orchestrate the creation of the initialization token.
	o.rw.Lock()

	// 4. Double-check: Essential because another goroutine might have
	// initialized o.once during the gap between RUnlock() and Lock().
	if o.once != nil {
		res = o.once
		o.rw.Unlock()
		select {
		case <-ctx.Done():
			var zero T
			return zero, true, ctx.Err()
		case <-o.opts.ctx.Done():
			var zero T
			return zero, true, o.opts.ctx.Err()
		case <-res.done:
			return res.value, true, res.err
		}
	}

	// 5. This specific goroutine wins the race and spawns the initialization task.
	res = &onceResult[T]{
		done: make(chan struct{}),
	}
	o.once = res

	// CRITICAL: Release the write lock immediately!
	// This unblocks all subsequent requests, allowing them to pass through the fast-path (RLock)
	// and safely block on res.done rather than queuing up outside o.rw.Lock().
	o.rw.Unlock()

	// 6. Launch an asynchronous goroutine to execute f.
	// This satisfies the requirement where the primary caller can return instantly upon
	// client-side context cancellation, while the background initialization continues undisturbed.
	go func() {
		// Execute f using the long-lived o.opts.ctx to guarantee completion.
		res.value, res.err = f(o.opts.ctx)

		// 7. Error handling & Cache Eviction.
		if res.err != nil {
			o.rw.Lock()
			// Only evict the cache if no other faster retry generation has overwritten o.once.
			if o.once == res {
				o.once = nil
			}
			o.rw.Unlock()
		}

		// Execute the custom callback if registered via WithDone.
		if o.opts.done != nil {
			o.opts.done(res.value, res.err)
		}

		// 8. Always broadcast the completion signal LAST.
		// This guarantees that any retry triggered by waiters waking up from close(res.done)
		// will deterministically see o.once as nil and start a fresh initialization generation.
		close(res.done)
	}()

	// 9. Await result or client-side cancellation.
	select {
	case <-ctx.Done():
		var zero T
		return zero, true, ctx.Err() // Client aborted early -> marked as shared=true to guard cache layers.
	case <-o.opts.ctx.Done():
		var zero T
		return zero, true, o.opts.ctx.Err()
	case <-res.done:
		// If the background task failed, treat shared as true so nobody writes the error to the upper cache.
		if res.err != nil {
			return res.value, true, res.err
		}
		return res.value, false, nil // The true winner! Returns shared = false.
	}
}
