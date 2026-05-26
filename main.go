package main

// import (
// 	"log"

// 	"github.com/os-webui/os-webui/cmd"
// )

// func main() {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)
// 	if e := cmd.Execute(); e != nil {
// 		log.Fatalln(e)
// 	}
// }
import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/os-webui/os-webui/internal/plugins"
)

func main() {
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

	// c, e := plugins.NewPlugin(`/home/dev/project/go/github/os-webui/bin/plugins/install`,
	// 	`fail2ban`)

	// fmt.Println(c, e)
	// i := interp.New(interp.Options{
	// 	SourcecodeFilesystem: plugins.FS(),
	// })
	// i.Use(stdlib.Symbols)

	// p, e := i.CompilePath(`fail2ban`)
	// if e != nil {
	// 	log.Fatalln(e)
	// }
	// fmt.Println(p.PackageName())
}
