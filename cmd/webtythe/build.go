// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/attic-labs/noms/go/d"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	prod := kingpin.Flag("prod", "builds for the production os/arch").Bool()

	_, dir, _, ok := runtime.Caller(0)
	d.PanicIfFalse(ok)
	dir = path.Dir(dir)
	ui := path.Join(dir, "ui")

	os.Setenv("GO111MODULE", "on")

	run(ui, "npm", "run", "build")
	run(ui, "go", "generate")

	c := []string{"go", "build"}
	if *prod {
		c = []string{"gox", "-osarch", "linux/amd64"}
	}

	run(dir, c[0], c[1:]...)
	fmt.Println("Done.")
}

func run(dir, c string, args ...string) {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	err := cmd.Run()
	d.PanicIfError(err)
}
