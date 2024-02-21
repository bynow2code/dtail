package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"testing"
)

func TestNewTailDir(t *testing.T) {
	dirname := "./testdata"
	NewTailDir(dirname)
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name string `json:"name"`
}

func TestDoUpdate(t *testing.T) {
	url := "https://api.github.com/repos/bynow2code/dtail/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("最新版本：", release.TagName)

	systemOS := runtime.GOOS
	if systemOS == "darwin" {
		systemOS = "macos"
	}
	systemARCH := runtime.GOARCH

	compressionFormat := ".tar.gz"

	filename := fmt.Sprintf("dtail_%s_%s_%s_%s", release.TagName, systemOS, systemARCH, compressionFormat)
	fmt.Println(filename)
}
