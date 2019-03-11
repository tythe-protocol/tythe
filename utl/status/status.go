// Package status prints status messages to a console, overwriting previous values.
package status

import (
	"fmt"
	"sync"
)

const (
	clearLine = "\x1b[2K\r"
)

var (
	m = sync.Mutex{}
)

func Clear() {
	m.Lock()
	fmt.Print(clearLine)
	m.Unlock()
}

func Enter() {
	m.Lock()
	fmt.Println()
	m.Unlock()
}

func Printf(format string, args ...interface{}) {
	m.Lock()
	fmt.Printf(clearLine+format, args...)
	m.Unlock()
}

type Writer struct{}

func (w Writer) Write(p []byte) (n int, err error) {
	if p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}
	Printf(string(p))
	return len(p), nil
}
