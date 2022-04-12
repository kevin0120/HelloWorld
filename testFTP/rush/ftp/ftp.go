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

type Plugin struct {
	Type     string            `yaml:"type"`
	Url      string            `yaml:"url"`
	Path     string            `yaml:"path"`
	Format   string            `yaml:"format"`
	User     string            `yaml:"user"`
	PassWord string            `yaml:"password"`
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers"`
}

var config = &Plugin{
	Type:     "ftp",
	Url:      "127.0.0.1:21",
	Path:     "rush/ftp",
	Format:   "",
	User:     "admin",
	PassWord: "admin",
	Method:   "",
	Headers:  map[string]string{},
}

var items = map[string]string{"1": "hello ftp", "2": "kevin"}

type Ftp struct {
	Url  string
	Path string

	username  string
	password  string
	ftpClient *ftp.ServerConn
}

func NewFtpProvider(config *Plugin) *Ftp {
	return &Ftp{
		Url:      config.Url,
		Path:     config.Path,
		username: config.User,
		password: config.PassWord,
	}
}

func (f *Ftp) changeDir() error {
	err := f.ftpClient.ChangeDir(f.Path)
	if err != nil {
		//这里不能这样写，不同的ftp服务器报的错不同
		if !strings.Contains(err.Error(), "found") {
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

func (f *Ftp) Write(data map[string]string) error {

	for k, item := range data {
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
	f := NewFtpProvider(config)
	_ = f.Connect()

	_ = f.Write()
	_ = f.Close()
}
