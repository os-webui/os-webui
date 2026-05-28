package web

import (
	"log/slog"
	"os"

	"github.com/os-webui/os-webui/config"
)

func Run(cfg *config.Config) error {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := initplugin(&cfg.Plugins, log)
	if err != nil {
		return nil
	}
	return runWeb(&cfg.Web, cfg.Dev, log)
}
