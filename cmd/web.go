package cmd

import (
	"log/slog"
	"os"

	"github.com/os-webui/os-webui/cmd/web"
	"github.com/os-webui/os-webui/config"
	"github.com/spf13/cobra"
)

const (
	defaultNetwork = `tcp`
	defaultAddr    = `:9026`
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
			cfg.Plugins.Install = "plugins"
			cfg.Plugins.Data = "plugins-data"

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
			web.Run(&cfg)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, `config`,
		`c`,
		`os-webui.yaml`,
		`config file`,
	)
	flags.StringVarP(&network, `network`,
		`N`,
		``,
		`listen network (default "`+defaultNetwork+`")`,
	)
	flags.StringVarP(&addr, `addr`,
		`A`,
		``,
		`listen address (default "`+defaultAddr+`")`,
	)
	rootCmd.AddCommand(cmd)
}
