package internal

import (
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
		t.Error(err)
	}

	err = file.DoUpgrade()
	if err != nil {
		t.Error(err)
	}

	//fmt.Printf("%#v \n", file.Name)
	//fmt.Printf("%#v \n", file.DownloadUrl)
	//fmt.Printf("%#v \n", file.LocalPath)

}
