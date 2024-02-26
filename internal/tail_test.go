package internal

import (
	"testing"
)

func TestNewTailDir(t *testing.T) {
	dirname := "./testdata"
	NewTailDir(dirname)
}
