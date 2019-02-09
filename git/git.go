package git

import (
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"
	gogit "gopkg.in/src-d/go-git.v4"
)

// Clone fetches the latest copy of the git repo at a URL to a local directory.
func Clone(url *url.URL, dataPath string) (rootPath string, err error) {
	w := func(err error) error {
		return errors.Wrap(err, "Could not clone/sync git repo")
	}
	// TODO: would be nice to use something more readable than sha1
	// couldn't find a handy path escaping function
	dirName := sha1.Sum([]byte(url.String()))
	fullPath := path.Join(dataPath, string(hex.EncodeToString(dirName[:])))
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
