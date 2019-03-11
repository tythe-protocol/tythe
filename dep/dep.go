// Package dep implements dependency crawling for tythe.
package dep

import (
	"fmt"

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

// ID is a unique identifier for a dependency.
type ID struct {
	// Type is the type of dependency.
	Type Type `json:"type"`

	// Name is a human-readable name for the dependency.
	Name string `json:"name"`
}

func (id ID) IsEmpty() bool {
	return id.Type == None && id.Name == ""
}

func (id ID) String() string {
	return fmt.Sprintf("%s:%s", id.Type, id.Name)
}

// Dep describes a dependency (direct or indirect) of a root package.
type Dep struct {
	ID

	// Conf is the Tythe config for the dependency, or nil if there is none.
	Conf *conf.Config `json:"config"`
}
