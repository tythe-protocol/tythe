package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/tythe-protocol/tythe/cmd/flags"
	"github.com/tythe-protocol/tythe/cmd/webtythe/ui"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("webtythe", "The tythe.dev server.")
	cacheDir := flags.CacheDir(app)
	port := app.Flag("port", "The port to serve on").Default("8080").Uint16()
	app.Parse(os.Args[1:])

	http.HandleFunc("/-/list", list(*cacheDir))
	http.HandleFunc("/-/ua", ua)
	http.Handle("/", http.FileServer(ui.Fs))
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
