package dep

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/go-tythe/conf"
	"github.com/tythe-protocol/go-tythe/git"
)

// Dep describes a dependency (direct or indirect) of a root package.
type Dep struct {
	// Name is a human-readable name for the dependency.
	Name string
	// Conf is the Tythe config for the dependency, or nil if there is none.
	Conf *conf.Config
}

// List returns the transitive dependencies of the module at <url>.
//
// If <url> has empty scheme or "file" scheme, it is interpreted as a path and read
// from there directly. Otherwise, it is downloaded to a temporary directly and read
// from the temp directory.
//
// Currently this only does Go dependencies, but it will grow to use many strategies.
func List(url *url.URL, cacheDir string) ([]Dep, error) {
	var p string
	var err error

	w := func(err error) error {
		return errors.Wrapf(err, "Cannot list dependencies of package: %s", url.String())
	}

	switch url.Scheme {
	case "":
		p = url.String()
		fmt.Println(os.Expand(p, nil))

		_, err := os.Stat(p)
		if err != nil {
			if os.IsNotExist(err) {
				err = fmt.Errorf("Directory does not exist: %s", p)
			}
			return nil, w(err)
		}
		break
	default:
		p, err = git.Clone(url, cacheDir)
		if err != nil {
			return nil, err
		}
		break
	}

	// It would be cool to use https://golang.org/src/cmd/go/internal/modload/
	// instead, but not allowed.
	err = os.Chdir(p)
	if err != nil {
		return nil, w(err)
	}

	defer os.Chdir("-")

	out, err := exec.Command("go", "list", "-f", "{{ .Path }}@{{ .Version }} {{ .Dir }}", "-m", "all").Output()
	if err != nil {
		return nil, w(err)
	}

	r := []Dep{}
	s := bufio.NewScanner(bytes.NewReader(out))

	// Skip the first line - it is the root module itself
	s.Scan()
	for s.Scan() {
		t := s.Text()
		p := strings.Split(t, " ")
		if len(p) != 2 {
			return nil, w(fmt.Errorf("Unexpected output from `go list`: %s", t))
		}
		name, dir := p[0], p[1]
		c, err := conf.Read(dir)
		if err != nil {
			return nil, w(err)
		}
		r = append(r, Dep{
			Name: name,
			Conf: c,
		})
	}

	return r, nil
}
