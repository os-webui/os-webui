package plugins

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type PluginInfo struct {
	ID       string   `yaml:"id" json:"id,omitempty"`
	Version  string   `yaml:"version" json:"version,omitempty"`
	Author   string   `yaml:"author" json:"author,omitempty"`
	Name     string   `yaml:"name" json:"name,omitempty"`
	Platform []string `yaml:"platform" json:"platform,omitempty"`

	Description string `yaml:"description" json:"description,omitempty"`
}

type PluginMeta struct {
	PluginInfo `yaml:",inline"`
	I18n       map[string]map[string]string `yaml:"i18n"`
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
