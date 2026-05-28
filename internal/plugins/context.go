package plugins

import (
	"context"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/gophini/coalesce"
	"github.com/os-webui/os-webui/sdk"
	bolt "go.etcd.io/bbolt"
)

type Context struct {
	*contextLow
	ui             sdk.WebUI
	ctx            context.Context
	acceptLanguage string
}

func (c *Context) UI() sdk.WebUI {
	return c.ui
}
func (c *Context) Context() context.Context {
	if c.ctx == nil {
		return c.contextLow.ctx
	}
	return c.ctx
}
func (c *Context) AcceptLanguage() string {
	return c.acceptLanguage
}

type contextLow struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	log *slog.Logger

	id, version, install, config, data string
	db                                 *coalesce.Task[*bolt.DB]

	keys map[string]any
	rw   sync.RWMutex
}

func newContext(log *slog.Logger, install, config, data, id, version string) (*contextLow, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &contextLow{
		ctx:        ctx,
		cancelFunc: cancelFunc,

		log: log,

		install: filepath.Join(install, id),
		config:  filepath.Join(config, id),
		data:    filepath.Join(data, id),

		id:      id,
		version: version,
		db:      coalesce.NewTask(coalesce.WithTaskContext[*bolt.DB](ctx)),
		keys:    make(map[string]any),
	}, nil
}
func (c *contextLow) Log() *slog.Logger {
	return c.log
}
func (c *contextLow) ID() string {
	return c.id
}
func (c *contextLow) Version() string {
	return c.version
}
func (c *contextLow) Data() string {
	return c.data
}
func (c *contextLow) Install() string {
	return c.install
}

func (c *contextLow) Config() string {
	return c.config
}

func (c *contextLow) DB(ctx context.Context) (*bolt.DB, error) {
	db, _, err := c.db.Get(ctx, func(ctx context.Context) (*bolt.DB, error) {
		return bolt.Open(filepath.Join(c.config, `plugin.db`), 0600, &bolt.Options{
			Timeout: time.Minute,
		})
	})
	return db, err
}
func (c *contextLow) Get(key string) (any, bool) {
	c.rw.RLock()
	val, ok := c.keys[key]
	c.rw.RUnlock()
	return val, ok
}

func (c *contextLow) Set(key string, val any) {
	c.rw.Lock()
	c.keys[key] = val
	c.rw.Unlock()
}

func (c *contextLow) Delete(key string) {
	c.rw.Lock()
	delete(c.keys, key)
	c.rw.Unlock()
}
