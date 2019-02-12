package dep

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/go-tythe/conf"
)

// Type represents a dependency "type". For example, Go or npm.
type Type int

const (
	// None represents no dependency type.
	None Type = iota

	// Go represents a golang dependency.
	Go
)

func (t Type) String() string {
	switch t {
	case None:
		return "<none>"
	case Go:
		return "go"
	default:
		return "<invalid>"
	}
}

// Dep describes a dependency (direct or indirect) of a root package.
type Dep struct {
	// Type is the type of dependency.
	Type Type

	// Name is a human-readable name for the dependency.
	Name string

	// Conf is the Tythe config for the dependency, or nil if there is none.
	Conf *conf.Config
}

func (d Dep) String() string {
	return fmt.Sprintf("%s:%s", d.Type, d.Name)
}

// List returns the transitive dependencies of the module at <path>.
//
// Currently this only does Go dependencies, but it will grow to use many strategies.
func List(path string) ([]Dep, error) {
	w := func(err error) error {
		return errors.Wrapf(err, "Cannot list dependencies of package: %s", path)
	}

	// It would be cool to use https://golang.org/src/cmd/go/internal/modload/
	// instead, but not allowed.
	err := os.Chdir(path)
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
			Type: Go,
			Name: name,
			Conf: c,
		})
	}

	return r, nil
}
