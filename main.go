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

	"github.com/os-webui/os-webui/internal/plugins"
)

func main() {
	c, e := plugins.NewPlugin(`/home/dev/project/go/github/os-webui/bin/plugins/install`,
		`fail2ban`)

	fmt.Println(c, e)
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
