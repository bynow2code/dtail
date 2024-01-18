package util

import (
	"fmt"
	"os"
)

func PrintFatalError(a ...any) {
	fmt.Fprintln(os.Stderr, a)
	os.Exit(1)
}
