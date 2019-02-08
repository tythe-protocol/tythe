package dep

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/go-tythe/git"
)

// List returns the transitive dependencies of the module at <url>.
// Currently this only does Go dependencies, but it will grow to use many strategies.
func List(url *url.URL) ([]string, error) {
	w := func(err error) error {
		return errors.Wrapf(err, "Could not list dependencies of package: %s", url.String())
	}
	dir, err := ioutil.TempDir("", "go-tythe")
	if err != nil {
		return nil, w(err)
	}

	p, err := git.Clone(url, dir)
	if err != nil {
		return nil, w(err)
	}

	// It would be cool to use https://golang.org/src/cmd/go/internal/modload/
	// instead, but not allowed.
	err = os.Chdir(p)
	if err != nil {
		return nil, w(err)
	}

	defer os.Chdir("-")

	out, err := exec.Command("go", "list", "-m", "all").Output()
	if err != nil {
		return nil, w(err)
	}

	r := []string{}
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		p := strings.Split(s.Text(), " ")[0]
		r = append(r, p)
	}

	// TODO: get the URL of the package, not the path

	return r, nil
}
