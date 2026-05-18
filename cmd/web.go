package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/os-webui/os-webui/cmd/web"
	"github.com/os-webui/os-webui/config"
	"github.com/spf13/cobra"
)

const (
	defaultNetwork = `tcp`
	defaultAddr    = `:9026`
)

var (
	defaultPluginsInstall = filepath.Join(`plugins`, `install`)
	defaultPluginsData    = filepath.Join(`plugins`, `data`)
	defaultPluginsConfig  = filepath.Join(`plugins`, `config`)
)

func init() {
	var (
		filename      string
		network, addr string
	)

	cmd := &cobra.Command{
		Use:   `web`,
		Short: `run web server`,
		Run: func(cmd *cobra.Command, args []string) {
			var cfg config.Config
			cfg.Web.Network = defaultNetwork
			cfg.Web.Addr = defaultAddr
			cfg.Plugins.Install = defaultPluginsInstall
			cfg.Plugins.Data = defaultPluginsData
			cfg.Plugins.Config = defaultPluginsConfig

			e := config.LoadConfig(filename, &cfg)
			if e != nil {
				slog.Error(e.Error())
				os.Exit(1)
			}
			if network != `` {
				cfg.Web.Network = network
			}
			if addr != `` {
				cfg.Web.Addr = addr
			}
			e = web.Run(&cfg)
			if e != nil {
				os.Exit(1)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, `config`,
		`c`,
		`os-webui.js`,
		`config file`,
	)
	flags.StringVarP(&network, `network`,
		`n`,
		``,
		`listen network (default "`+defaultNetwork+`")`,
	)
	flags.StringVarP(&addr, `addr`,
		`a`,
		``,
		`listen address (default "`+defaultAddr+`")`,
	)
	rootCmd.AddCommand(cmd)
}
