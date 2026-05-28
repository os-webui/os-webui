package plugins

import (
	"log/slog"
	"strings"
	"sync"
)

var DefaultPluginsManager = NewPluginsManager()

type PluginsManager struct {
	plugins map[string]*Plugin
	rw      sync.RWMutex
}

func NewPluginsManager() *PluginsManager {
	return &PluginsManager{
		plugins: map[string]*Plugin{},
	}
}
func (m *PluginsManager) Add(id string, plugin *Plugin) {
	m.rw.Lock()
	m.plugins[id] = plugin
	m.rw.Unlock()
}
func (m *PluginsManager) Plugin(id string) *Plugin {
	m.rw.Lock()
	plugin := m.plugins[id]
	m.rw.Unlock()
	return plugin
}

func (m *PluginsManager) Cleanup(slog *slog.Logger) {
	var waits []<-chan struct{}
	for _, plugin := range m.plugins {
		metadata := plugin.Metadata()
		slog.Info(`plugin cleanup`, `id`, metadata.ID, `version`, metadata.Version)
		ch := plugin.plugin.Quit()
		if ch != nil {
			waits = append(waits, ch)
		}
		plugin.Cleanup()
	}
	for _, ch := range waits {
		<-ch
	}
}
func (m *PluginsManager) List(acceptLanguage string) []PluginInfo {
	strs := strings.FieldsFunc(acceptLanguage, func(r rune) bool {
		return r == ',' || r == ';' || r == ' '
	})
	langs := make([]string, 0, len(strs))
	for _, s := range strs {
		s = strings.TrimSpace(s)
		if s == `` || strings.Contains(s, `=`) {
			continue
		}
		langs = append(langs, s)
	}
	m.rw.RLock()
	defer m.rw.RUnlock()
	items := make([]PluginInfo, 0, len(m.plugins))
	for _, p := range m.plugins {
		metadata := p.metadata
		info := metadata.PluginInfo
		i18n := metadata.I18n
		if len(i18n) != 0 {
			for _, lang := range langs {
				if keys, ok := i18n[lang]; ok {
					if name, ok := keys[`name`]; ok && name != `` {
						info.Name = name
					}
					if description, ok := keys[`description`]; ok && description != `` {
						info.Description = description
					}
					break
				}
				i := strings.Index(lang, `-`)
				if i > 0 {
					lang = lang[:i]
					if keys, ok := i18n[lang]; ok {
						if name, ok := keys[`name`]; ok && name != `` {
							info.Name = name
						}
						if description, ok := keys[`description`]; ok && description != `` {
							info.Description = description
						}
						break
					}
				}
			}
		}

		items = append(items, info)
	}
	return items
}
func (m *PluginsManager) Get(id, acceptLanguage string) (PluginInfo, bool) {
	strs := strings.FieldsFunc(acceptLanguage, func(r rune) bool {
		return r == ',' || r == ';' || r == ' '
	})
	langs := make([]string, 0, len(strs))
	for _, s := range strs {
		s = strings.TrimSpace(s)
		if s == `` || strings.Contains(s, `=`) {
			continue
		}
		langs = append(langs, s)
	}
	m.rw.RLock()
	defer m.rw.RUnlock()
	p, ok := m.plugins[id]
	if !ok {
		return PluginInfo{}, false
	}

	metadata := p.metadata
	info := metadata.PluginInfo
	i18n := metadata.I18n
	if len(i18n) != 0 {
		for _, lang := range langs {
			if keys, ok := i18n[lang]; ok {
				if name, ok := keys[`name`]; ok && name != `` {
					info.Name = name
				}
				if description, ok := keys[`description`]; ok && description != `` {
					info.Description = description
				}
				break
			}
			i := strings.Index(lang, `-`)
			if i > 0 {
				lang = lang[:i]
				if keys, ok := i18n[lang]; ok {
					if name, ok := keys[`name`]; ok && name != `` {
						info.Name = name
					}
					if description, ok := keys[`description`]; ok && description != `` {
						info.Description = description
					}
					break
				}
			}
		}
	}
	return info, true
}
