package env

import (
	"fmt"
	"os"
)

// Must gets the environment variable with the specified name or panics if it doesn't exist.
func Must(name string) string {
	v := os.Getenv(name)
	if v == "" {
		fmt.Fprintf(os.Stderr, "Could not find required environment variable: %s\n", name)
		os.Exit(1)
	}
	return v
}
