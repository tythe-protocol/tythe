package golang

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tythe-protocol/tythe/conf"
	"github.com/tythe-protocol/tythe/dep/shared"

	"github.com/pkg/errors"
)

func w(err error, path string) error {
	return errors.Wrapf(err, "Cannot list go dependencies of package: %s", path)
}

// List returns all the transitive Go dependencies of the package at <path>
func List(path string) (r []shared.Dep, err error) {
	// It would be cool to use https://golang.org/src/cmd/go/internal/modload/
	// instead, but not allowed.
	err = os.Chdir(path)
	if err != nil {
		return nil, w(err, path)
	}

	defer os.Chdir("-")

	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		return nil, w(err, path)
	}

	out, err := exec.Command("go", "list", "-f", "{{ .Path }}@{{ .Version }} {{ .Dir }}", "-m", "all").Output()
	if err != nil {
		return nil, w(err, path)
	}

	s := bufio.NewScanner(bytes.NewReader(out))

	// Skip the first line - it is the root module itself
	s.Scan()
	for s.Scan() {
		t := s.Text()
		p := strings.Split(t, " ")
		if len(p) != 2 {
			return nil, w(fmt.Errorf("Unexpected output from `go list`: %s", t), path)
		}
		name, dir := p[0], p[1]
		c, err := conf.Read(dir)
		if err != nil {
			return nil, w(err, path)
		}
		r = append(r, shared.Dep{
			Name: name,
			Conf: c,
		})
	}

	return r, nil
}
