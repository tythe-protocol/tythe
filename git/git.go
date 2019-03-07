package git

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"
	gogit "gopkg.in/src-d/go-git.v4"
)

// DirForURL calculates a unique filename to store a repo in.
func DirForURL(url *url.URL, dataDir string) string {
	// TODO: would be nice to use something more readable than sha1
	// couldn't find a handy path escaping function
	dirName := sha1.Sum([]byte(url.String()))
	return path.Join(dataDir, string(hex.EncodeToString(dirName[:])))
}

// Clone fetches the latest copy of the git repo at a URL to a local directory.
func Clone(url *url.URL, dataPath string) (rootPath string, err error) {
	w := func(err error) error {
		return errors.Wrap(err, "Could not clone/sync git repo")
	}

	fullPath := DirForURL(url, dataPath)
	_, err = os.Stat(fullPath)

	if err == nil {
		// repo already exists, sync it.
		var wt *gogit.Worktree
		r, err := gogit.PlainOpen(fullPath)
		if err == nil {
			wt, err = r.Worktree()
			if err == nil {
				err = wt.Pull(&(gogit.PullOptions{}))
			}
		}
		if err != nil && err != gogit.NoErrAlreadyUpToDate {
			return "", w(err)
		}
		return fullPath, nil
	}

	if !os.IsNotExist(err) {
		return "", w(err)
	}

	_, err = gogit.PlainClone(fullPath, false, &gogit.CloneOptions{
		URL:   url.String(),
		Depth: 1,
	})
	return fullPath, err
}

// Resolve resolves a URL to a local or remote git repository to one that is local.
func Resolve(url *url.URL, cacheDir string) (path string, err error) {
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
		path, err = Clone(url, cacheDir)
		if err != nil {
			return "", err
		}
	}
	return
}
