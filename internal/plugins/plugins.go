package plugins

import (
	"log/slog"
)

var DefaultPluginsManager = NewPluginsManager()

type PluginsManager struct {
	plugins map[string]*Plugin
}

func NewPluginsManager() *PluginsManager {
	return &PluginsManager{
		plugins: map[string]*Plugin{},
	}
}
func (m *PluginsManager) Add(id string, plugin *Plugin) {
	m.plugins[id] = plugin
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
