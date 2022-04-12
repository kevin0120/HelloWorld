// author: smallfish <smallfish.xy@gmail.com>

package ftp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type FTP struct {
	host    string
	port    int
	user    string
	passwd  string
	pasv    int
	cmd     string
	Code    int
	Message string
	Debug   bool
	stream  []byte
	conn    net.Conn
	Error   error
}

func (ftp *FTP) debugInfo(s string) {
	if ftp.Debug {
		fmt.Println(s)
	}
}

func (ftp *FTP) Connect(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	if ftp.conn, ftp.Error = net.Dial("tcp", addr); ftp.Error != nil {
		return
	}
	ftp.Response()
	ftp.host = host
	ftp.port = port
}

func (ftp *FTP) Login(user, passwd string) {
	ftp.Request("USER " + user)
	ftp.Request("PASS " + passwd)
	ftp.user = user
	ftp.passwd = passwd
}

func (ftp *FTP) Response() (code int, message string) {
	ret := make([]byte, 1024)
	n, _ := ftp.conn.Read(ret)
	msg := string(ret[:n])
	code, _ = strconv.Atoi(msg[:3])
	message = msg[4 : len(msg)-2]
	ftp.debugInfo("<*cmd*> " + ftp.cmd)
	ftp.debugInfo(fmt.Sprintf("<*code*> %d", code))
	ftp.debugInfo("<*message*> " + message)
	return
}

func (ftp *FTP) Request(cmd string) {
	ftp.conn.Write([]byte(cmd + "\r\n"))
	ftp.cmd = cmd
	ftp.Code, ftp.Message = ftp.Response()
	//此处是因为 如果这个命令错误 比如说Mkd命令错误 它会发两条回复
	if ftp.Code == 0 {
		ftp.Code, ftp.Message = ftp.Response()
	}

	if cmd == "EPSV" {
		start, end := strings.Index(ftp.Message, "(|||"), strings.Index(ftp.Message, "|)")
		if start < 0 {
			ftp.pasv = 0
			return
		}
		s := strings.Split(ftp.Message[start+4:end], ",")
		l1, _ := strconv.Atoi(s[0])
		ftp.pasv = l1
		fmt.Printf("@@@@@@@@@@@@@@@@@@\n%s\n%d@@@@@@@@@@@@@@\n", ftp.Message, ftp.pasv)
	}
	if (cmd != "EPSV") && (ftp.pasv > 0) {
		ftp.Message = newRequest(ftp.host, ftp.pasv, ftp.stream)
		ftp.debugInfo("<*response*> " + ftp.Message)
		ftp.pasv = 0
		ftp.stream = nil
		ftp.Code, _ = ftp.Response()
	}
}

func (ftp *FTP) Pasv() {
	ftp.Request("EPSV")
}

func (ftp *FTP) Pwd() {
	ftp.Request("PWD")
}

func (ftp *FTP) Cwd(path string) {
	ftp.Request("CWD " + path)
}

func (ftp *FTP) Mkd(path string) {
	ftp.Request("MKD " + path)
}

func (ftp *FTP) Size(path string) (size int) {
	ftp.Request("SIZE " + path)
	size, _ = strconv.Atoi(ftp.Message)
	return
}

func (ftp *FTP) List() {
	ftp.Pasv()
	ftp.Request("LIST")
}

func (ftp *FTP) Stor(file string, data []byte) {
	ftp.Pasv()
	if data != nil {
		ftp.stream = data
	}
	ftp.Request("STOR " + file)
}

func (ftp *FTP) Quit() {
	ftp.Request("QUIT")
	ftp.conn.Close()
}

// new connect to FTP pasv port, return data
func newRequest(host string, port int, b []byte) string {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	defer conn.Close()
	if b != nil {
		conn.Write(b)
		return "OK"
	}
	ret := make([]byte, 4096)
	n, _ := conn.Read(ret)
	return string(ret[:n])
}
