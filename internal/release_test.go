package internal

import (
	"fmt"
	"testing"
)

func TestGithubRelease_Latest(t *testing.T) {
	release := NewGithubRelease()
	release.Latest()
	fmt.Printf("%#v \n", release)
	fmt.Printf("%#v \n", release.UpgradeFile().Name)
}
