// Package status prints status messages to a console, overwriting previous values.
// All the functions in this package (including Writer.Write) are safe to be called
// from concurrent goroutines.
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

// Printf works like fmt.Printf, except that it overwrites the current line first.
func Printf(format string, args ...interface{}) {
	m.Lock()
	fmt.Printf(clearLine+format, args...)
	m.Unlock()
}

// Clear clears the current line.
func Clear() {
	m.Lock()
	fmt.Print(clearLine)
	m.Unlock()
}

// Enter moves to the next line.
func Enter() {
	m.Lock()
	fmt.Println()
	m.Unlock()
}

// Writer implements io.Writer by sending to status.Printf.
type Writer struct{}

func (w Writer) Write(p []byte) (n int, err error) {
	if p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}
	Printf(string(p))
	return len(p), nil
}
