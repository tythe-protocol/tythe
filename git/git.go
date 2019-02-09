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
	// TODO: would be nice to use something more readable than sha1
	// couldn't find a handy path escaping function
	dirName := sha1.Sum([]byte(url.String()))
	fullPath := path.Join(dataPath, string(hex.EncodeToString(dirName[:])))
	_, err = os.Stat(fullPath)

	if err != nil && !os.IsNotExist(err) {
		return "", errors.Wrapf(err, "Could not clone into directory: %s", fullPath)
	}

	if err == nil {
		// TODO: sync
		return fullPath, nil
	}

	_, err = gogit.PlainClone(fullPath, false, &gogit.CloneOptions{
		URL:   url.String(),
		Depth: 1,
	})
	return fullPath, err
}
