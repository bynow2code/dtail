/*
Copyright © 2024 changqq <https://github.com/bynow2code/dtail>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"bufio"
	"dTail/cmd"
	"dTail/util"
	"fmt"
	"github.com/inconshreveable/go-update"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("检测到新版本，是否现在升级(y/n): ")
	reader := bufio.NewReader(os.Stdin)
	readString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	text := strings.Replace(readString, "\n", "", -1)
	if text == "y" {
		fmt.Println("开始升级...")

		url := "https://github.com/bynow2code/dtail/releases/download/v0.0.3/dtail_0.0.3_macos_arm64.tar.gz"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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

	cmd.Execute()
}
