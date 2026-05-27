package main

import (
	"log"

	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/os-webui/os-webui/cmd"

	"github.com/os-webui/os-webui/internal/plugins"
)

func main() {
	RunMain()
	// RunPlugin()
}
func RunMain() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}

func RunPlugin() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	install, e := filepath.Abs(`plugins/install`)
	if e != nil {
		return
	}
	config, e := filepath.Abs(`plugins/config`)
	if e != nil {
		return
	}
	data, e := filepath.Abs(`plugins/data`)
	if e != nil {
		return
	}
	name := `fail2ban`
	plugin, e := plugins.New(log, install, config, data, name)
	if e != nil {
		return
	}
	fmt.Println(plugin)
}
