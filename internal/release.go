package internal

import (
	"encoding/json"
	"fmt"
	"github.com/bynow2code/dtail/util"
	"io"
	"net/http"
	"runtime"
)

var Version string

type Release interface {
	Latest()
	Version() string
	UpgradeFile() *UpgradeFile
}

type UpgradeFile struct {
	Name        string
	DownloadUrl string
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

func (g *GithubRelease) Latest() {
	url := "https://api.github.com/repos/bynow2code/dtail/releases/latest"
	response, err := http.Get(url)
	if err != nil {
		util.PrintError("请求Github发生错误", response)
	}
	jsonStr, err := io.ReadAll(response.Body)
	if err != nil {
		util.PrintError("读取Github返回值发生错误", err)
	}
	err = json.Unmarshal(jsonStr, g)
	if err != nil {
		util.PrintError("解析Github返回值", err)
	}
}

func (g *GithubRelease) UpgradeFile() *UpgradeFile {
	return &UpgradeFile{
		Name:        upgradeFileName(g.Version()),
		DownloadUrl: "",
	}
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
