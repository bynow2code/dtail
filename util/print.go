package util

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func PrintlnErrorf(format string, a ...any) {
	PrintlnError(fmt.Sprintf(format, a...))
}

func PrintlnError(a any) {
	_, _ = color.New(color.FgRed).Fprintf(os.Stderr, "Error: %v\n", a)
}

func PrintlnFatalf(format string, a ...any) {
	PrintlnFatal(fmt.Sprintf(format, a...))
}

func PrintlnFatal(a any) {
	_, _ = color.New(color.FgRed).Add(color.Bold).Fprintf(os.Stderr, "Fatal error: %v\n", a)
	os.Exit(1)
}

func PrintlnInfof(format string, a ...any) {
	PrintlnInfo(fmt.Sprintf(format, a...))
}

func PrintlnInfo(a any) {
	_, _ = color.New(color.FgBlue).Printf("Info: %v\n", a)
}
