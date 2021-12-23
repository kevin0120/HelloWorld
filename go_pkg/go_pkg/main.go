package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
)

func handleError(err error)  {
	fmt.Println(err)
}

func main() {
	fr, err := os.Open("dmeo.tar") // 打开tar包文件，返回*io.Reader
	handleError(err)               // handleError为错误处理函数，下同
	defer fr.Close()

	// 实例化新的tar.Reader
	tr := tar.NewReader(fr)

	for {
		hdr, err := tr.Next() // 获取下一个文件，第一个文件也用此方法获取
		if err == io.EOF {
			break // 已读到文件尾
		}
		handleError(err)

		// 通过创建文件获得*io.Writer
		fw, _ := os.Create("demo/" + hdr.Name)
		handleError(err)

		// 拷贝数据
		_, err = io.Copy(fw, tr)
		handleError(err)
	}
}
