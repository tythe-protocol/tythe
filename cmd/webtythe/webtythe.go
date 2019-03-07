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
	app.Parse(os.Args[1:])

	http.HandleFunc("/-/list", list(*cacheDir))
	http.HandleFunc("/-/ua", ua)
	http.Handle("/", http.FileServer(ui.Fs))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
