package internal

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/bynow2code/dtail/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var Version string

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

func NewRelease() *Release {
	return &Release{}
}

func (r *Release) Latest() {
	util.PrintInfo("开始读取最新release")
	url := "https://api.github.com/repos/bynow2code/dtail/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		util.PrintError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.PrintError(err)
	}

	if err = json.Unmarshal(body, r); err != nil {
		util.PrintError(err)
	}
	util.PrintInfo("最新版本：", r.TagName)

	goos := runtime.GOOS
	if goos == "darwin" {
		goos = "macos"
	}
	goarch := runtime.GOARCH
	fileName := fmt.Sprintf("dtail_%s_%s_%s.tar.gz", r.TagName, goos, goarch)
	util.PrintInfo("准备获取文件下载链接：", fileName)

	var downloadUrl string
	for _, assetInfo := range r.Assets {
		if assetInfo.Name == fileName {
			downloadUrl = assetInfo.DownloadUrl
			break
		}
	}
	util.PrintInfo("下载地址：", downloadUrl)

	util.PrintInfo("开始下载")
	get, err := http.Get(url)
	if err != nil {
		util.PrintError(err)
	}
	util.PrintInfo("下载完成")

	unzip(get.Body)
}

func unzip(r io.Reader) {
	fmt.Println("开始解压")
	gr, err := gzip.NewReader(r)
	if err != nil {
		util.PrintError(err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			util.PrintError(err)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			util.PrintError(err)
		}

		curFile := filepath.Join(home, hdr.Name)
		dir := filepath.Dir(curFile)
		if _, err = os.Open(dir); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(filepath.Dir(curFile), 0755)
				if err != nil {
					util.PrintError(err)
				}
			} else {
				util.PrintError(err)
			}
		}

		err = os.MkdirAll(filepath.Dir(curFile), 0755)
		if err != nil {
			util.PrintError(err)
		}

		fo, err := os.Create(curFile)
		if err != nil {
			util.PrintError(err)
		}
		defer fo.Close()

		if _, err := io.Copy(fo, tr); err != nil {
			util.PrintError(err)
		}

		util.PrintInfo(hdr.Linkname)
	}
	util.PrintInfo("解压完成")

	//util.PrintInfo("准备升级")
	//open, err := os.Open("/Users/edy/dtail_0.0.3_macos_arm64/dtail")
	//if err != nil {
	//	return
	//}
	//err = update.Apply(open, update.Options{})
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//util.PrintInfo("升级完成")
	//os.Exit(0)
}
