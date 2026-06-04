package sdk

import (
	"context"
	"log/slog"

	bolt "go.etcd.io/bbolt"
)

// UI defines a reactive presentation layer interface for dynamic data exchange
// between the system backend and front-end layout rendering engines.
type WebUI interface {
	// Get extracts an input argument mapped to the specific layout identifier.
	// It decodes the payload into the destination object and returns its raw representation.
	Get(id string, jsonObjectPointer any) bool

	// Put serializes and flushes the given rendering value back into the designated UI structure node.
	Put(id string, jsonObject any)
}

// Context provides a highly optimized, thread-safe localized environment runtime execution context.
// It fully wraps the native Go context.Context, empowering pluggable extensions with lifecycle-aware control streams.
type Context interface {
	Context() context.Context
	AcceptLanguage() string

	Log() *slog.Logger

	// ID returns the unique canonical registration identifier of the running plugin.
	ID() string

	// Version returns the current semantic versioning string of the loaded plugin instance.
	Version() string

	// Install returns the read-only file system directory path where the plugin's source/assets reside.
	Install() string

	// Config returns the path targeting the stateful customized configuration file directory assigned to this plugin.
	Config() string

	// Data returns the path targeting the stateful runtime operational asset/log directory assigned to this plugin.
	Data() string

	// Get retrieves a volatile, transient runtime variable from the in-memory execution context storage track.
	Get(key string) (any, bool)

	// Set mutates or registers a transient runtime variable inside the memory storage scope.
	Set(key string, val any)

	// Delete atomitally purges a specified transient runtime variable from the active context window.
	Delete(key string)

	// DB compiles or attaches a high-performance transactional embedded B+ tree engine (bbolt) instance.
	// This yields zero-network-overhead persistence dedicated strictly to the caller plugin's storage space.
	// Data() + `plugin.db`
	DB(context.Context) (*bolt.DB, error)

	// 讀取配置，如果 name 是相對路徑，會 加載 filepath.Join(Config(),name)
	LoadConfig(ctx context.Context, name string) ([]byte, error)
	// 儲存配置，如果 name 是相對路徑，會 儲存到 filepath.Join(Config(),name)
	// 會先 Save 的 name+`.temp` 再 Rename 到 name 以防止寫入截斷數據
	SaveConfig(ctx context.Context, name string, data []byte) error

	// UI provisions access to the interactive user interface rendering pipeline coupled to the active execution cycle.
	UI() WebUI
}

// Plugin specifies the standardized structural lifecycle control loop contracts required to construct
// secure, hot-swappable dynamically decoupled extension programs within the server management panel.
type Plugin interface {
	// Quit returns a read-only synchronization signal channel.
	// It closes abruptly when a fatal lifecycle exception forces the extension instance to drop and terminate.
	Quit() <-chan struct{}

	// OnStartup acts as the core entry-point hook triggered right after the dynamic bytecode compiler finishes structural linking.
	// External resource allocations, background tasks, and daemon loops should be initialized here.
	OnStartup(ctx Context) error

	// OnCleanup operates as the terminal lifecycle destructor hook invoked during graceful panel shutdowns
	// or dynamic plug-and-play hot-unloading sequences to clean volatile allocations and file locks.
	OnCleanup(ctx Context)

	// 加載配置文件供網頁修改，如果不需要特殊處理通常應該調用 ctx.LoadConf
	LoadConfig(ctx Context, name string) (string, error)
	// 網頁回調保存配置文件，如果不需要特殊處理通常應該直接調用 ctx.SaveConf
	SaveConfig(ctx Context, name string, data string) error

	// 插件應該重新載入配置文件
	OnReload(ctx Context) error

	// Features inventories a list of declarative capability blocks exposed by this extension to the host orchestration engine.
	Features(ctx Context) []Feature
}

// Feature represents an atomic, isolated operational routing capability block exposed to front-end UI clients.
type Feature interface {
	ID() string
	// Metadata outputs a standardized JSON Schema or declarative schema template string.
	// Front-end dashboard components utilize this schema blueprint to render automated input/output forms and validate contracts.
	Metadata(ctx Context) (M, error)

	// Run executes the core business capability logic of the target feature.
	// Inputs and structural responses are multiplexed implicitly via the provided responsive sdk.Context model.
	Run(ctx Context) (err error)
}
type M map[string]any
