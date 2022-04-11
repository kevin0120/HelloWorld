package main

import (
	"bytes"
	"github.com/jlaffaye/ftp"
	"strings"
	"time"
)

const (
	FtpDialWithTimeout = 5 * time.Second
)

var items = map[string]string{"1": "hello ftp", "2": "kevin"}

type Ftp struct {
	Url  string
	Path string

	username  string
	password  string
	ftpClient *ftp.ServerConn
}

func NewFtpProvider() *Ftp {
	return &Ftp{
		Url:      "127.0.0.1:21",
		Path:     "rush/client",
		username: "admin",
		password: "admin",
	}
}

func (f *Ftp) changeDir() error {
	err := f.ftpClient.ChangeDir(f.Path)
	if err != nil {
		//这里不能这样写，不同的ftp服务器报的错不同
		if !strings.Contains(err.Error(), "not found") {
			err = f.ftpClient.MakeDir(f.Path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (f *Ftp) Connect() error {
	c, err := ftp.Dial(f.Url, ftp.DialWithTimeout(FtpDialWithTimeout))
	if err != nil {
		return err
	}

	err = c.Login(f.username, f.password)
	if err != nil {
		return err
	}
	f.ftpClient = c
	err = f.changeDir()
	if err != nil {
		return err
	}
	return nil
}

func (f *Ftp) Write() error {

	for k, item := range items {
		data := bytes.NewBuffer([]byte(item))
		err := f.ftpClient.Stor(k, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Ftp) Close() error {
	if err := f.ftpClient.Quit(); err != nil {
		return err
	}
	return nil
}

func main() {
	f := NewFtpProvider()
	_ = f.Connect()

	_ = f.Write()
	_ = f.Close()
}
