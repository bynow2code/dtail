package internal

import (
	"fmt"
	"testing"
)

func TestGithubRelease_Latest(t *testing.T) {
	release := NewGithubRelease()
	err := release.Latest()
	if err != nil {
		t.Error(err)
	}

	file, err := release.UpgradeFile()
	if err != nil {
		return
	}

	fmt.Printf("%#v \n", file.Name)
	fmt.Printf("%#v \n", file.DownloadUrl)
}
