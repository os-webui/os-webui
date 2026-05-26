package lazy

import (
	"context"
	"sync"
)

type Cache[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Add(key K, value V)
}

// loaderOptions bundles the generic context and backend storage interface.
type loaderOptions[K comparable, V any] struct {
	ctx   context.Context
	cache Cache[K, V]
}

// Loader manages concurrent read/write coalescing across a keyspace.
// It leverages localized Value[V] generations to isolate hot-key synchronization.
type Loader[K comparable, V any] struct {
	opts loaderOptions[K, V]
	rw   sync.RWMutex

	fetchers map[K]*Value[V]
}

// WithLoaderContext attaches a long-lived supervisor context to the cache broker.
func WithLoaderContext[K comparable, V any](ctx context.Context) Option[*loaderOptions[K, V]] {
	return newFuncOption(func(o *loaderOptions[K, V]) {
		if ctx == nil {
			o.ctx = context.Background()
		} else {
			o.ctx = ctx
		}
	})
}

// WithLoaderCache configures the secondary persistence or eviction tier (e.g., an LRU cache).
func WithLoaderCache[K comparable, V any](cache Cache[K, V]) Option[*loaderOptions[K, V]] {
	return newFuncOption(func(o *loaderOptions[K, V]) {
		o.cache = cache
	})
}

// NewLoader initializes a ready-to-use cluster-safe loader manager.
func NewLoader[K comparable, V any](opt ...Option[*loaderOptions[K, V]]) *Loader[K, V] {
	loader := &Loader[K, V]{
		opts: loaderOptions[K, V]{
			ctx: context.Background(),
		},
		fetchers: make(map[K]*Value[V]), // FIX: Initialized map to prevent nil-pointer mutation panic.
	}
	for _, o := range opt {
		o.apply(&loader.opts)
	}
	return loader
}

// Get returns (value, shared, error).
// It utilizes a strict hierarchical orchestration pipeline: Map State Check -> LRU Settled Check -> Flight Node Creation.
func (c *Loader[K, V]) Get(ctx context.Context, key K, fetcher func(context.Context) (V, error)) (V, bool, error) {
	// 1. Fast-Path (RLock): Intercept hot queries if the computation is in flight or already settled.
	c.rw.RLock()
	if loader, ok := c.fetchers[key]; ok {
		c.rw.RUnlock()
		return loader.Get(ctx, fetcher)
	}
	if c.opts.cache != nil {
		if value, found := c.opts.cache.Get(key); found {
			c.rw.RUnlock()
			return value, true, nil
		}
	}
	c.rw.RUnlock()

	// 2. Slow-Path (Lock): Acquire mutation lock for synchronization gatekeeper setup.
	c.rw.Lock()
	if loader, ok := c.fetchers[key]; ok {
		c.rw.Unlock()
		return loader.Get(ctx, fetcher)
	}
	if c.opts.cache != nil {
		if value, found := c.opts.cache.Get(key); found {
			c.rw.Unlock()
			return value, true, nil
		}
	}

	// 3. Orchestration: Hook the cache persistence and map eviction atomically into WithDone.
	// This generic callback executes exactly once when the underlying execution routine resolves.
	loader := New(WithContext[V](c.opts.ctx), WithDone(func(val V, err error) {
		c.rw.Lock()
		delete(c.fetchers, key)
		if err == nil && c.opts.cache != nil {
			c.opts.cache.Add(key, val) // Perfectly safe and guarantees atomic handover
		}
		c.rw.Unlock()
	}))

	c.fetchers[key] = loader
	c.rw.Unlock()

	// 4. Delegation: Hand over execution context to the distinct in-flight Value agent.
	return loader.Get(ctx, fetcher)
}
