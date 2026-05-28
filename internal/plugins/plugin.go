package plugins

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"runtime"

	"github.com/os-webui/os-webui/internal/symbols"
	"github.com/os-webui/os-webui/sdk"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

// 🌟 Regular Expression Explanation:
// ^[a-zA-Z]      : Assures the string starts with an uppercase or lowercase English letter.
// [a-zA-Z0-9_-]* : Followed by zero or more letters, numbers, hyphens, or underscores.
// $              : Assures the match extends to the absolute end of the string.
var matchID = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// MatchID validates if the given plugin or feature ID adheres to the panel's naming convention.
func MatchID(id string) bool {
	return matchID.MatchString(id)
}

type Plugin struct {
	plugin   sdk.Plugin
	metadata *PluginMeta
	ctx      *Context
}

func New(log *slog.Logger,
	install, config, data,
	id string) (ret *Plugin, e error) {
	log = log.With(`plugin`, id)
	metadata, e := LoadPluginMeta(install, id)
	if e != nil {
		log.Error(`failed to load plugin meta`, `error`, e)
		return
	}
	if len(metadata.Platform) != 0 {
		platform := runtime.GOOS + `/` + runtime.GOARCH
		ok := false
		for _, v := range metadata.Platform {
			if platform == v {
				ok = true
				break
			}
		}
		if !ok {
			log.Warn(platform+` is not on the list of supported platforms`, `platform`, metadata.Platform)
			return
		}
	}
	plugin, e := NewPlugin(install, id)
	if e != nil {
		log.Error(`failed to new plugin`, `error`, e)
		return
	}

	cl, e := newContext(log, install, config, data, id, metadata.Version)
	if e != nil {
		log.Error(`failed to new Context`, `error`, e)
		return
	}
	ctx := &Context{
		contextLow: cl,
	}
	ret = &Plugin{
		plugin:   plugin,
		metadata: metadata,
		ctx:      ctx,
	}
	return
}

func NewPlugin(dir, name string) (ret sdk.Plugin, err error) {
	i := interp.New(interp.Options{
		GoPath:               magicDir,
		SourcecodeFilesystem: FS(dir),
		Unrestricted:         true,
	})
	e := i.Use(stdlib.Symbols)
	if e != nil {
		return nil, e
	}
	if e = i.Use(syscall.Symbols); e != nil {
		return nil, e
	}
	if e = os.Setenv("YAEGI_SYSCALL", "1"); e != nil {
		return nil, e
	}
	if e = i.Use(unsafe.Symbols); e != nil {
		return nil, e
	}
	if e = os.Setenv("YAEGI_UNSAFE", "1"); e != nil {
		return nil, e
	}
	if e = i.Use(unrestricted.Symbols); e != nil {
		return nil, e
	}
	if e = os.Setenv("YAEGI_UNRESTRICTED", "1"); e != nil {
		return nil, e
	}
	e = i.Use(symbols.Symbols)
	if e != nil {
		return nil, e
	}

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf(`%v`, p)
		}
	}()
	// Eval
	_, e = i.EvalPath(name)
	if e != nil {
		return nil, e
	}
	pkgs := i.Symbols(name)
	if pkgs == nil {
		return nil, errors.New(`not found Symbols: ` + name)
	}

	keys := pkgs[name]
	if keys == nil {
		return nil, errors.New(`not found Package: ` + name)
	}
	fi, ok := keys[`New`]
	if !ok {
		return nil, errors.New(`not found func: New`)
	}
	pluginNew, ok := fi.Interface().(func() sdk.Plugin)
	if !ok {
		return nil, errors.New(`func signature mismatch, should func() (sdk.Plugin)`)
	}
	return pluginNew(), nil
}
func (p *Plugin) Metadata() *PluginMeta {
	return p.metadata
}
func (p *Plugin) Startup() error {
	return p.plugin.OnStartup(p.ctx)
}

func (p *Plugin) Cleanup() {
	p.plugin.OnCleanup(p.ctx)
}

type FeatureInfo struct {
	ID       string `json:"id,omitempty"`
	Metadata sdk.M  `json:"metadata,omitempty"`
}

func (p *Plugin) Features(c context.Context, acceptLanguage string) ([]FeatureInfo, error) {
	ctx, cancel := context.WithCancelCause(c)
	defer cancel(context.Canceled)
	go func() {
		select {
		case <-p.ctx.contextLow.ctx.Done():
			cancel(p.ctx.contextLow.ctx.Err())
		case <-ctx.Done():
			cancel(ctx.Err())
		}
	}()

	newCTX := &Context{
		contextLow:     p.ctx.contextLow,
		ctx:            ctx,
		acceptLanguage: acceptLanguage,
	}
	features := p.plugin.Features(newCTX)
	items := make([]FeatureInfo, len(features))
	for i, feature := range features {
		id := feature.ID()
		metadata, err := feature.Metadata(newCTX)
		if err != nil {
			return nil, err
		}
		items[i] = FeatureInfo{
			ID:       id,
			Metadata: metadata,
		}
	}
	return items, nil
}
