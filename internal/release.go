package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var Version string

var (
	ErrDownloadUrlEmpty = errors.New("download url is empty")
)

type Release interface {
	Latest() error
	Version() string
	UpgradeFile() (*UpgradeFile, error)
}

type UpgradeFile struct {
	Name         string
	DownloadUrl  string
	DownloadPath string
}

type GithubRelease struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

func NewGithubRelease() Release {
	return &GithubRelease{}
}

func (g *GithubRelease) Version() string {
	return g.TagName
}

func (g *GithubRelease) Latest() error {
	url := "https://api.github.com/repos/bynow2code/dtail/releases/latest"
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	jsonStr, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonStr, g)
	if err != nil {
		return err
	}
	return nil
}

func (g *GithubRelease) UpgradeFile() (*UpgradeFile, error) {
	upgrade := &UpgradeFile{}
	upgrade.Name = upgradeFileName(g.Version())
	for _, asset := range g.Assets {
		if asset.Name == upgrade.Name {
			upgrade.DownloadUrl = asset.DownloadUrl
			break
		}
	}
	if upgrade.DownloadUrl == "" {
		return nil, ErrDownloadUrlEmpty
	}

	return upgrade, nil
}

func upgradeFileName(version string) string {
	goos := runtime.GOOS
	if goos == "darwin" {
		goos = "macos"
	}
	arch := runtime.GOARCH
	//dtail_v0.0.1-alpha.1_linux_amd64.tar.gz
	filename := fmt.Sprintf("dtail_%s_%s_%s.tar.gz", version, goos, arch)
	return filename
}

type Upgrade interface {
	Download() error
	Extract() error
}

func (f *UpgradeFile) Download() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	downloadPath := filepath.Join(dir, f.Name)
	open, err := os.OpenFile(downloadPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	response, err := http.Get(f.DownloadUrl)
	if err != nil {
		return err
	}
	_, err = io.Copy(open, response.Body)
	if err != nil {
		return err
	}
	f.DownloadPath = downloadPath
	return nil
}

func (f *UpgradeFile) Extract() error {
	return nil
}
