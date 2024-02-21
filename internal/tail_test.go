package internal

import (
	"dTail/util"
	"net/http"
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

	//err := update.Apply(resp.Body, update.Options{})
	//if err != nil {
	//	// error handling
	//}
	//return err
}
