package internal

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/inconshreveable/go-update"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var Version string

const (
	UpgradeFileName = "dtail"
)

var (
	ErrDownloadUrlEmpty = errors.New("download url is empty")
)

type Release interface {
	Latest() error
	Version() string
	UpgradeFile() (Upgrade, error)
}

type Upgrade interface {
	DoUpgrade() error
}

type GzTarUpgradeFile struct {
	Name        string
	DownloadUrl string
	LocalPath   string
	UpgradePath string
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
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

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

func (g *GithubRelease) UpgradeFile() (Upgrade, error) {
	upgrade := &GzTarUpgradeFile{}
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

func (f *GzTarUpgradeFile) download() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	f.LocalPath = filepath.Join(dir, f.Name)

	open, err := os.OpenFile(f.LocalPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer open.Close()

	client := http.Client{
		Timeout: 60 * time.Second,
	}
	response, err := client.Get(f.DownloadUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(open, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (f *GzTarUpgradeFile) extract() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	open, err := os.Open(f.LocalPath)
	if err != nil {
		return err
	}
	defer open.Close()

	gzr, err := gzip.NewReader(open)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			targetFile := filepath.Join(dir, header.Name)
			targetPath := filepath.Dir(targetFile)
			err = os.MkdirAll(targetPath, 0755)
			if err != nil {
				return err
			}

			create, err := os.Create(targetFile)
			if err != nil {
				return err
			}
			defer create.Close()

			_, err = io.Copy(create, tr)
			if err != nil {
				return err
			}

			targetBase := filepath.Base(targetFile)
			if targetBase == UpgradeFileName {
				f.UpgradePath = targetFile
			}
		}
	}
	return nil
}

func (f *GzTarUpgradeFile) DoUpgrade() error {
	err := f.download()
	if err != nil {
		return err
	}

	err = f.extract()
	if err != nil {
		return err
	}

	open, err := os.Open(f.UpgradePath)
	if err != nil {
		return err
	}

	err = update.Apply(open, update.Options{})
	return err
}
