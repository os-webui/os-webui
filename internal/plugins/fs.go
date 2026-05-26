package plugins

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const magicDir = `48f8dc68666c4283a8887e0c25303c0b_5182346e341446cb8f4943f254226fbc`

var magicDirSrc = filepath.Join(magicDir, `src`)

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
