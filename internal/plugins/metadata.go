package plugins

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type PluginMeta struct {
	ID       string   `yaml:"id"`
	Version  string   `yaml:"version"`
	Author   string   `yaml:"author"`
	Name     string   `yaml:"name"`
	Platform []string `yaml:"platform"`

	Description string `yaml:"description"`

	I18n map[string]map[string]string `yaml:"i18n"`
}

func LoadPluginMeta(dir, id string) (ret *PluginMeta, e error) {
	data, e := os.ReadFile(filepath.Join(dir, id, `plugin.yaml`))
	if e != nil {
		return
	}
	var metadata PluginMeta
	e = yaml.Unmarshal(data, &metadata)
	if e != nil {
		return
	}
	if metadata.ID != id {
		e = errors.New(`plugin id not matched: ` + metadata.ID)
		return
	}
	ret = &metadata
	return
}
