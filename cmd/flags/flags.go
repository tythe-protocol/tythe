package flags

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type flaggable interface {
	Flag(name string, description string) *kingpin.FlagClause
}

func CacheDir(cmd flaggable) *string {
	hd, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return cmd.Flag("cache-dir", "Directory to write cached repos to during crawling").
		Default(fmt.Sprintf("%s/.tythe", hd)).String()
}

func Sandbox(cmd flaggable) *bool {
	return cmd.Flag("sandbox", "Use the sandbox Coinbase API").Default("false").Bool()
}
