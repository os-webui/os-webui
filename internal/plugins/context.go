package plugins

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"github.com/os-webui/os-webui/internal/lazy"
	bolt "go.etcd.io/bbolt"
)

type Context struct {
	*contextLow
}
type contextLow struct {
	context.Context
	keys                      *sync.Map
	cancelFunc                context.CancelFunc
	install, config, data, id string
	db                        *lazy.Value[*bolt.DB]
}

func newContext(install, config, data, id string) (*contextLow, error) {

	ctx, cancelFunc := context.WithCancel(context.Background())
	return &contextLow{
		Context:    ctx,
		cancelFunc: cancelFunc,
		install:    filepath.Join(install, id),
		config:     filepath.Join(config, id),
		data:       filepath.Join(data, id),
		id:         id,
		keys:       &sync.Map{},
		db:         lazy.New[*bolt.DB](lazy.WithContext[*bolt.DB](ctx)),
	}, nil
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
func (c *contextLow) ID() string {
	return c.id
}
func (c *contextLow) Map() *sync.Map {
	return c.keys
}

// 	// Get input
// 	Input(id string) (any, bool)
// 	// Get input
// 	InputString(id string) (string, bool)

// 	// Push output
// 	OutputString(id string, value string)

func (c *contextLow) DB(ctx context.Context) (*bolt.DB, error) {
	db, _, err := c.db.Get(ctx, func(ctx context.Context) (*bolt.DB, error) {
		return bolt.Open(filepath.Join(c.config, `plugin.db`), 0600, &bolt.Options{
			Timeout: time.Minute,
		})
	})
	return db, err
}
