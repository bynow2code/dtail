package util

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func PrintlnErrorf(format string, a ...any) {
	PrintlnError(fmt.Sprintf(format, a...))
}

func PrintlnError(a ...any) {
	_, _ = color.New(color.FgRed).Fprintln(os.Stderr, a...)
}

func PrintlnFatalf(format string, a ...any) {
	PrintlnFatal(fmt.Sprintf(format, a...))
}

func PrintlnFatal(a ...any) {
	_, _ = color.New(color.FgRed).Add(color.Bold).Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func PrintlnInfof(format string, a ...any) {
	PrintlnInfo(fmt.Sprintf(format, a...))
}

func PrintlnInfo(a ...any) {
	_, _ = color.New(color.FgBlue).Println(a...)
}
