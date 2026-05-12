package logger

import (
	"fmt"
	"os"
)

// Error levels
const (
	FATAL  = "fatal"
	ACCESS = "access"
)

var hadError bool

func PrintError(level, path string, err error) {
	switch level {
	case FATAL:
		fmt.Fprintf(os.Stderr, "my-ls: %v\n", err)
		fmt.Fprintln(os.Stderr, "Try 'my-ls --help' for more information.")
		os.Exit(2)

	case ACCESS:
		hadError = true
		fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': %v\n", path, err)
	}
}

func ExitStatus() {
	if hadError {
		os.Exit(2)
	}
}
