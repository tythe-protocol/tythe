package git

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// DirForURL calculates a unique filename to store a repo in.
func DirForURL(url *url.URL, dataDir string) string {
	// TODO: would be nice to use something more readable than sha1
	// couldn't find a handy path escaping function
	dirName := sha1.Sum([]byte(url.String()))
	return path.Join(dataDir, string(hex.EncodeToString(dirName[:])))
}

func runGit(wd string, l *log.Logger, args ...string) error {
	var err error

	for i := 0; ; i++ {
		l.Printf("running %+v in %s\n", args, wd)

		// For mysterious reasons, Git sometimes fails to exit.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Dir = wd
		err = cmd.Run()

		if err == nil {
			return nil
		}

		if i == 2 {
			return err
		}

		l.Println("Retrying...")
	}
}

// Clone fetches the latest copy of the git repo at a URL to a local directory.
func Clone(url *url.URL, dataPath string, l *log.Logger) (rootPath string, err error) {
	fullPath := DirForURL(url, dataPath)
	_, err = os.Stat(fullPath)

	if err == nil {
		// repo already exists, sync it.
		err := runGit(fullPath, l, "git", "fetch", "--depth", "1")
		if err != nil {
			return "", errors.Wrapf(err, "Could not pull: %s: %s", url.String(), err.Error())
		}
		return fullPath, nil
	}

	if !os.IsNotExist(err) {
		return "", errors.Wrapf(err, "Could not stat: %s: %s", url.String(), err.Error())
	}

	if url.Hostname() == "github.com" {
		// TODO: Which protocol is fastest? ssh, https, or git?
		url.Scheme = "https"
	}
	err = runGit("", l, "git", "clone", "--depth", "1", url.String(), fullPath)
	if err != nil {
		return "", errors.Wrapf(err, "Could not clone: %s: %s", url.String(), err.Error())
	}
	return fullPath, err
}

// Resolve resolves a URL to a local or remote git repository to one that is local.
func Resolve(url *url.URL, cacheDir string, l *log.Logger) (path string, err error) {
	if url.Scheme == "" {
		path = url.String()
		path, err = filepath.Abs(path)
		if err != nil {
			return "", errors.Wrapf(err, "Could not get absolute path of %s: %s", path, err.Error())
		}
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				err = fmt.Errorf("Directory does not exist: %s", path)
			}
			return "", errors.Wrapf(err, "Could not resolve package: %s", url.String())
		}
		// TODO: Probably want to copy to a temporary directory to avoid modifying target dir.
	} else {
		path, err = Clone(url, cacheDir, l)
		if err != nil {
			return "", err
		}
	}
	return
}
