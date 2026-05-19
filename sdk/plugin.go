package sdk

import (
	"context"

	bolt "go.etcd.io/bbolt"
)

// Context provides localized environment isolation parameters for a running plugin
type Context interface {
	context.Context

	// Config returns the stateful customized configuration file directory path
	Config() string

	// Data returns the stateful database/log storage file directory path
	Data() string

	// Locale returns the standard client language identifier, e.g., 'en', 'zh-Hant'
	Locale() string

	// Get input
	Input(id string) (any, bool)
	// Get input
	InputString(id string) (string, bool)

	// Push output
	OutputString(id string, value string)

	Set(id string, value []byte) error
	Get(id string) []byte
	DB() (*bolt.DB, error)
}

// Plugin defines the standardized runtime control loop lifecycle for external extensions
type Plugin interface {
	// Quit returns a read-only channel emitting severe exceptions causing plugin termination
	Quit() <-chan struct{}

	// OnStartup triggers immediately after the dynamic script engine finishes initial loading
	OnStartup(ctx Context) error

	// OnCleanup triggers during systems shutdown or plugin dynamic hot-unloading routines
	OnCleanup(ctx Context)

	// Features inventories all declarative capability operational blocks exposed by this plugin
	Features(ctx Context) []Feature
}

// Feature represents a single functional block executing actions within the system panel
type Feature interface {
	// Metadata exposes structural inputs, outputs maps for UI rendering engines
	Metadata(ctx Context) (string, error)

	// Run executes the feature command with input payloads and yields structured responses or error states
	Run(ctx Context) (err error)
}
