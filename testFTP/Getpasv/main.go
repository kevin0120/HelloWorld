package main

import (
	"fmt"
	"github.com/kevin0120/HelloWorld/testFTP/Getpasv/Rewrite"
	"io/ioutil"
	"os"
)

func main() {
	ftp := new(ftp.FTP)
	// debug default false
	ftp.Debug = true
	ftp.Connect("localhost", 21)

	// login
	ftp.Login("admin", "admin")
	if ftp.Code == 530 {
		fmt.Println("error: login failure")
		os.Exit(-1)
	}

	// pwd
	ftp.Pwd()

	fmt.Println("code:", ftp.Code, ", message:", ftp.Message)

	// make dir
	//ftp.Mkd("/path")
	//ftp.Request("TYPE I")
	//ftp.Request("EPSV")
	//ftp.Request("EPSV")
	//
	//ftp.Request("EPSV")
	//ftp.Request("EPSV")
	ftp.List()

	//ftp.Request("EPSV")
	// stor file
	b, _ := ioutil.ReadFile("./compile")
	ftp.Stor("/path/a.txt", b)

	ftp.Quit()
}
