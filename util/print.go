package util

import (
	"fmt"
	"os"
)

func PrintError(a ...any) {
	fmt.Fprint(os.Stderr, "dtail:")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func PrintInfo(a ...any) {
	fmt.Fprint(os.Stdout, "dtail:")
	fmt.Fprintln(os.Stdout, a...)
}
