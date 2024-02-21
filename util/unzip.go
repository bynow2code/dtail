package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Unzip(r io.Reader) {
	fmt.Println("开始解压")

	gr, err := gzip.NewReader(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	// 循环读取文件
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // 到达文件尾部，结束循环
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			return
		}

		curFile := filepath.Join(home, hdr.Name)
		err = os.MkdirAll(filepath.Dir(curFile), 0755)
		if err != nil {
			return
		}

		fo, err := os.Create(curFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fo.Close()

		if _, err := io.Copy(fo, tr); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(curFile)
	}

	fmt.Println("解压完成")
}
