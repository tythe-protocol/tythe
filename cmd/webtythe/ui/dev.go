// +build !release

package ui

import (
	"net/http"
	"os"
)

type nullfs struct{}

func (nfs nullfs) Open(name string) (http.File, error) {
	return nil, os.ErrNotExist
}

var Fs = nullfs{}
