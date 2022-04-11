example 
========
```bash

package main

import (
	"fmt"
	"github.com/smallfish/ftp"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

func Dirlist(conn net.Conn) (a []string) {
	ret := make([]byte, 1024)
	n, _ := conn.Read(ret)
	msg := string(ret[:n])
	//fmt.Println(msg)
	s := strings.Split(msg, "\n")
	//fmt.Println(len(s))
	for a := 0; a < len(s)-1; a++ {
		fmt.Println("@@@@@@@@@@@")
		fmt.Println(s[a])
		fmt.Println("###########")
		//err := ioutil.WriteFile(path, d, 0666)
	}
	return s
}

func Copyfile(conn net.Conn, filename string) (a []string) {
	ret := make([]byte, 1024*1024*10)
	n, _ := conn.Read(ret)
	msg := string(ret[:n])
	fmt.Println("@@@@@@@@@@@")
	//fmt.Println(msg)
	b := strings.Split(filename, "/")
	fmt.Println("/home/kevin/mytest/" + b[len(b)-1])
	err := ioutil.WriteFile("/home/kevin/mytest/"+b[len(b)-1], []byte(msg), 0666)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("###########")
	s := strings.Split(msg, "\n")

	return s
}

func main() {
	ftp := new(ftp.FTP)
	// debug default false
	ftp.Debug = true
	ftp.Connect("localhost", 21)

	// login
	ftp.Login("kevin", "123456")
	if ftp.Code == 530 {
		fmt.Println("error: login failure")
		os.Exit(-1)
	}

	// pwd
	ftp.Pwd()
	fmt.Println("code:", ftp.Code, ", message:", ftp.Message)

	// make dir
	//	ftp.Mkd("./path")
	ftp.Request("TYPE I")
	/*
		var ch chan []string
		var s []string
		go Dirlist("localhost",ftp.Getpasv(),ch)
		<-ch
	*/
	for {
		ftp.Request("PASV")
		fmt.Println(ftp.Getpasv())
		addr1 := fmt.Sprintf("%s:%d", "localhost", ftp.Getpasv())
		conn1, _ := net.Dial("tcp", addr1)
		ftp.Request("NLST /home/kevin/ftp")
		s := Dirlist(conn1)

		for a := 0; a < len(s)-1; a++ {
			ftp.Request("PASV")
			fmt.Println(ftp.Getpasv())
			addr2 := fmt.Sprintf("%s:%d", "localhost", ftp.Getpasv())
			conn2, _ := net.Dial("tcp", addr2)
			ftp.Request("RETR " + s[a])

			Copyfile(conn2, s[a])
			ftp.Request("DELE " + s[a])

		}
		time.Sleep(20 * time.Second)
	}
	//Copyfile("localhost",ftp.Getpasv())
	//ftp.Request("MKD /home/kevin/ftp/TEST1")

	//ftp.Request("PASV")
	//	  ftp.Request("LIST /home/kevin/ftp")

	//	ftp.Request("PASV")
	//	fmt.Println(ftp.Getpasv())

	//	ftp.Request("RMD /home/kevin/ftp/TEST1")

	//ftp.Request("CWD /home/kevin/ftp")
	//ftp.Pwd()

	// stor file
	//b, _ := ioutil.ReadFile("/home/kevin/changan/12.txt")
	//fmt.Println(string(b))

	//ftp.Stor("/home/uftp/12.txt", b)

	ftp.Quit()
}
```