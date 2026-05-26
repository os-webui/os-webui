package plugins

import (
	"errors"
	"log/slog"
	"os"

	"github.com/os-webui/os-webui/internal/symbols"
	"github.com/os-webui/os-webui/sdk"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

type Plugin struct {
	sdk.Plugin
	metadata *PluginMeta
	ctx      Context
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
	plugin, e := NewPlugin(install, id)
	if e != nil {
		log.Error(`failed to new plugin`, `error`, e)
		return
	}

	cl, e := newContext(install, config, data, id)
	if e != nil {
		log.Error(`failed to new Context`, `error`, e)
		return
	}

	ctx := Context{
		contextLow: cl,
	}
	ret = &Plugin{
		Plugin:   plugin,
		metadata: metadata,
		ctx:      ctx,
	}
	return
}

func NewPlugin(dir, name string) (sdk.Plugin, error) {
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
