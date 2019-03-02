package dep

import (
	"fmt"

	"github.com/tythe-protocol/tythe/dep/golang"
	"github.com/tythe-protocol/tythe/dep/npm"

	"github.com/tythe-protocol/tythe/conf"
)

// Type represents a dependency "type". For example, Go or npm.
type Type int

const (
	// None represents no dependency type.
	None Type = iota

	// Go represents a golang dependency.
	Go

	// NPM represents a npm dependency.
	NPM
)

func (t Type) String() string {
	switch t {
	case None:
		return "<none>"
	case Go:
		return "go"
	case NPM:
		return "npm"
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
	r := []Dep{}

	gd, err := golang.List(path)
	if err != nil {
		return nil, err
	}
	for _, d := range gd {
		r = append(r, Dep{
			Type: Go,
			Name: d.Name,
			Conf: d.Conf,
		})
	}

	nd, err := npm.List(path)
	if err != nil {
		return nil, err
	}
	for _, d := range nd {
		r = append(r, Dep{
			Type: NPM,
			Name: d.Name,
			Conf: d.Conf,
		})
	}

	return r, nil
}
