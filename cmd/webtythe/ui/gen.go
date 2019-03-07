// +build ignore

package main

import (
	"net/http"
	"path"
	"runtime"

	"github.com/attic-labs/noms/go/d"
	"github.com/shurcooL/vfsgen"
)

func main() {
	_, fn, _, _ := runtime.Caller(0)
	dir := http.Dir(path.Join(path.Dir(fn), "build"))
	err := vfsgen.Generate(dir, vfsgen.Options{
		PackageName:  "ui",
		BuildTags:    "release",
		VariableName: "Fs",
	})
	d.PanicIfError(err)
}
