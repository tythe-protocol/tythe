package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/go-tythe/git"
)

func resolvePackage(url *url.URL, cacheDir string) (path string, err error) {
	if url.Scheme == "" {
		path = url.String()
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				err = fmt.Errorf("Directory does not exist: %s", path)
			}
			return "", errors.Wrapf(err, "Could not resolve package: %s", url.String())
		}
	} else {
		path, err = git.Clone(url, cacheDir)
		if err != nil {
			return "", err
		}
	}
	return
}
