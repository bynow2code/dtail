package internal

import (
	"dTail/util"
	"fmt"
	"github.com/inconshreveable/go-update"
	"net/http"
	"os"
	"testing"
)

func TestNewTailDir(t *testing.T) {
	dirname := "./testdata"
	NewTailDir(dirname)
}

func TestDoUpdate(t *testing.T) {
	url := "https://github.com/bynow2code/dtail/releases/download/v0.0.3/dtail_0.0.3_macos_arm64.tar.gz"
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	util.Unzip(resp.Body)

	open, err := os.Open("/Users/edy/dtail_0.0.3_macos_arm64/dtail")
	if err != nil {
		return
	}
	err = update.Apply(open, update.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
