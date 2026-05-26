# Lazy - A Type-Safe, High-Concurrency Protective Utility Suite for Go

`lazy` is a generic, production-grade Go library designed to orchestrate complex concurrent execution lifecycles. It provides powerful primitives to solve classic high-concurrency pain points: **thundering herd problems (cache stampede), redundant lazy initialization, transient client cancellation leakage, and dynamic background telemetry tracking.**

Unlike rigid singleflight implementations, `lazy` empowers you to bind long-lived contexts to background workers, decouple caller timeouts from core tasks, gracefully evict failed states for retry logic, and automatically cancel background tasks if all callers abandon the request.

---

## Key Features

- **100% Type-Safe**: Built on Go Generics—no more frustrating `interface{}` (`any`) type assertions.
- **Thundering Herd Shield**: Consolidates thousands of concurrent identical requests into a single flight.
- **Graceful Failure Eviction**: Automatically evicts failed iterations, allowing subsequent requests to seamlessly trigger safe retries.
- **Atomic Handover Isolation**: Isolates the in-flight synchronization map from the final storage cache (e.g., LRU), eliminating race conditions during execution handovers.
- **Dynamic Orphan Cancellation**: Optional telemetry tracking via an atomic reference counter. If all callers timeout or abort, the background routine is terminated early to save CPU/IO.
- **Deadlock-Free Architectural Footprint**: Lock scopes are strictly localized to pure-memory mutations; third-party callbacks and cache storage additions run safely outside internal lock regions.

---

## Installation

```bash
go get [github.com/your-username/lazy](https://github.com/your-username/lazy)

```

---

## Components & Usage Examples

### 1. `lazy.Value[T]` (Single Resource Orchestration)

Perfect for lazily initializing heavy global singletons (e.g., database connection pools, plugin configurations) with resilient timeout/retry dynamics.

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"[github.com/your-username/lazy](https://github.com/your-username/lazy)"
	_ "[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)"
)

func main() {
	// Initialize a lazy value bound to the application lifetime context
	appCtx := context.Background()
	
	dbLoader := lazy.New[*sql.DB](
		lazy.WithContext[*sql.DB](appCtx),
		lazy.WithAutoCancel[*sql.DB](true), // Cancel DB connection attempt if all callers give up
		lazy.WithDone(func(db *sql.DB, err error) {
			if err != nil {
				fmt.Printf("Telemetry: Initialization failed: %v\n", err)
			} else {
				fmt.Println("Telemetry: Database pool successfully bootstrapped!")
			}
		}),
	)

	// Transient client request context with a short timeout
	clientCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// High-concurrency entry point
	db, shared, err := dbLoader.Get(clientCtx, func(bgCtx context.Context) (*sql.DB, error) {
		// This heavy operation uses bgCtx, ensuring it finishes even if clientCtx times out,
		// UNLESS all other interested waiters also timeout (thanks to WithAutoCancel).
		return sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hello")
	})

	if err != nil {
		panic(err)
	}

	// shared == false means this specific goroutine won the initialization race.
	// shared == true means it safely shared/awaited the result of the winner.
	fmt.Printf("Connected! Shared: %t, DB: %v\n", shared, db)
}

```

### 2. `lazy.Loader[K, V]` (Keyspaced Coalescing & Double-Cache Guard)

Perfect for protecting key-value caches (like LRU caches, Redis, or local maps) from cache stampedes. It creates temporary `Value[T]` nodes dynamically per key and safely populates your persistent cache exactly once.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"[github.com/your-username/lazy](https://github.com/your-username/lazy)"
)

// MockCache implements lazy.Cache using a simple map (or wrap your favorite LRU Cache here)
type MockCache struct {
	data map[string]string
}

func (m *MockCache) Get(key string) (string, bool) {
	v, ok := m.data[key]; return v, ok
}
func (m *MockCache) Add(key string, value string) {
	m.data[key] = value
}

func main() {
	underlyingLRU := &MockCache{data: make(map[string]string)}
	
	// Create a keyspaced Loader
	loader := lazy.NewLoader[string, string](
		lazy.WithLoaderContext[string, string](context.Background()),
		lazy.WithLoaderCache[string, string](underlyingLRU),
	)

	ctx := context.Background()
	
	// Simulate 10,000 concurrent requests querying the exact same expired key simultaneously
	// Only 1 database fetcher will execute. The rest will block at the channel layer and 
	// receive the populated cache item atomically without invoking data races on the LRU cache.
	val, hit, err := loader.Get(ctx, "user_123", func(bgCtx context.Context) (string, error) {
		time.Sleep(500 * time.Millisecond) // Simulate slow DB query
		return "John Doe", nil
	})

	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Result: %s, Cache Hit: %t\n", val, hit)
}

```

---

## Structural Execution Pipeline

The library orchestrates a seamless handoff pipeline when fetching elements through `Loader[K, V]` to safeguard state transitions:

1. **Fast-Path Check (`RLock`)**: Instantly yields the data if it is already present in your secondary Cache tier (e.g. LRU) or registers the request under an active in-flight worker.
2. **Slow-Path Barrier (`Lock`)**: If the key is entirely new, a unique localized synchronization generation node (`Value[T]`) is instantiated inside the tracking buffer.
3. **Atomic `WithDone` Handover**: When the background worker successfully completes the lookup task, it updates the secondary Cache and purges itself from the in-flight tracker under a unified memory fence. Subsequent lookups immediately bounce to the Cache, ensuring **zero duplicate write friction**.

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.
