package util

import (
	"fmt"
	"os"
)

func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

func PrintFatal(err error) {
	fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
	os.Exit(1)
}

func PrintInfo(a ...any) {
	fmt.Fprintln(os.Stdout, a...)
}
