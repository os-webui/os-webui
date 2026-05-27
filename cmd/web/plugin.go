package web

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/os-webui/os-webui/config"
	"github.com/os-webui/os-webui/internal/plugins"
)

func abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Abs(path)
}
func createPluginDir(tag, path string, slog *slog.Logger) (string, error) {
	dir, err := abs(path)
	if err != nil {
		slog.Error(`failed to abs plugin `+tag+` dir`, `dir`, path)
		return ``, err
	}
	slog.Info(`plugin `+tag+` dir`, `dir`, dir)
	err = os.MkdirAll(dir, 0644)
	if err != nil {
		slog.Error(`failed to mkdir plugin `+tag+` dir`, `dir`, dir)
		return ``, err
	}
	return dir, nil
}
func initplugin(cfg *config.PluginConfig, slog *slog.Logger) error {
	install, err := createPluginDir(`install`, cfg.Install, slog)
	if err != nil {
		return err
	}
	config, err := createPluginDir(`config`, cfg.Config, slog)
	if err != nil {
		return err
	}
	data, err := createPluginDir(`data`, cfg.Data, slog)
	if err != nil {
		return err
	}
	dirs, err := os.ReadDir(install)
	if err != nil {
		slog.Error(`failed to range plugin dirs`)
		return err
	}
	var items []*plugins.Plugin
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		id := dir.Name()
		if !plugins.MatchID(id) {
			continue
		}
		plugin, err := plugins.New(slog, install, config, data, id)
		if err != nil || plugin == nil {
			continue
		}
		items = append(items, plugin)
	}

	for _, plugin := range items {
		metadata := plugin.Metadata()
		slog.Info(`plugin startup`, `id`, metadata.ID, `version`, metadata.Version)
		err := plugin.Startup()
		if err != nil {
			slog.Warn(`failed to plugin startup`, `id`, metadata.ID, `version`, metadata.Version, `error`, err)
		}
		plugins.DefaultPluginsManager.Add(metadata.ID, plugin)
	}
	return nil
}
