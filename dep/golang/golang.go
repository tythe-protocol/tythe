package golang

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/tythe/dep"
)

func w(err error, path string) error {
	return errors.Wrapf(err, "Cannot list go dependencies of package: %s", path)
}

func Dir(name, dataDir string) string {
	return path.Join(build.Default.GOPATH, "pkg", "mod", name)
}

// Dependencies returns all the direct Golang dependencies from the module named <path>
func Dependencies(path string) ([]dep.ID, error) {
	// It would be cool to use https://golang.org/src/cmd/go/internal/modload/
	// instead, but not allowed.
	err := os.Chdir(path)
	if err != nil {
		return nil, w(err, path)
	}

	defer os.Chdir("-")

	// TODO: If there is a go.sum file we should just use it, rather than run tidy again.
	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		return nil, w(err, path)
	}

	out, err := exec.Command("go", "list", "-f", "{{ .Path }}@{{ .Version }} {{ .Dir }} {{ .Indirect }}", "-m", "all").Output()
	if err != nil {
		return nil, w(err, path)
	}

	s := bufio.NewScanner(bytes.NewReader(out))
	r := []dep.ID{}

	// Skip the first line - it is the root module itself
	s.Scan()
	for s.Scan() {
		t := s.Text()
		p := strings.Split(t, " ")
		if len(p) != 3 {
			return nil, w(fmt.Errorf("Unexpected output from `go list`: %s", t), path)
		}
		pathver, _, indirect := p[0], p[1], p[2]
		if indirect == "false" {
			r = append(r, dep.ID{
				Type: dep.Go,
				Name: pathver,
			})
		}
	}

	return r, nil
}
