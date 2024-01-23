package util

import (
	"fmt"
	"os"
)

func FatalError(a ...any) {
	fmt.Fprint(os.Stderr, "dtail：")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
