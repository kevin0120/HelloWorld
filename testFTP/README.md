FTP client for Go(lang)
==================================

install 
========
go get github.com/smallfish/ftp

example 
========
```bash
package main

import (
	"fmt"
	"github.com/smallfish/ftp"
	"io/ioutil"
	"os"
)

func main() {
	ftp := new(ftp.FTP)
	// debug default false
	ftp.Debug = true
	ftp.Connect("localhost", 21)

	// login
	ftp.Login("anonymous", "")
	if ftp.Code == 530 {
		fmt.Println("error: login failure")
		os.Exit(-1)
	}

	// pwd
	ftp.Pwd()
	fmt.Println("code:", ftp.Code, ", message:", ftp.Message)

	// make dir
	ftp.Mkd("/path")
	ftp.Request("TYPE I")

	// stor file
	b, _ := ioutil.ReadFile("/path/a.txt")
	ftp.Stor("/path/a.txt", b)

	ftp.Quit()
}
```


ftp主动模式被动模式 
```bash
https://www.cnblogs.com/trfizeng/p/4354832.html

ftp 服务器搭建
https://blog.csdn.net/weixin_44961081/article/details/118683045
```
