package plugins

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/os-webui/os-webui/internal/symbols"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

type FS string

func (s FS) Open(name string) (fs.File, error) {
	if !strings.HasPrefix(name, magicDirSrc) {
		return nil, os.ErrNotExist
	}
	name = filepath.Clean(filepath.Join(string(s), name[len(magicDirSrc):]))
	if len(name) < len(s)+1 ||
		name[len(s)] != filepath.Separator ||
		!strings.HasPrefix(name, string(s)) {
		return nil, os.ErrNotExist
	}
	return os.Open(name)
}

type Plugin struct {
}

const magicDir = `48f8dc68666c4283a8887e0c25303c0b_5182346e341446cb8f4943f254226fbc`

var magicDirSrc = filepath.Join(magicDir, `src`)

func NewPlugin(dir, name string) (*Plugin, error) {
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

	fmt.Println(keys)
	// if pkgs == nil {
	// 	return fmt.Errorf(`func () not found: %s`, name)
	// }
	// keys := pkgs[c.path]
	// if keys == nil {
	// 	return fmt.Errorf(`func () not found: %s`, name)
	// }
	return &Plugin{}, nil
}
