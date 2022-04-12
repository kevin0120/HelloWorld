package main


import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"
)

var items = map[string]string{"H1": "hello ftp", "H2": "kevin"}

type Sftp struct {
	Url  string
	Path string

	username   string
	password   string
	sftpClient *sftp.Client
}

func NewSftpProvider() *Sftp {
	return &Sftp{
		Url:      "192.168.1.7:22",
		Path:     "./",
		username: "kevin",
		password: "123456",
	}
}

func connect(user, password, addr string) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func (s *Sftp) pathExists() (bool, error) {
	_, err := s.sftpClient.Stat(s.Path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *Sftp) Connect() error {
	var ok bool
	sftpClient, err := connect(s.username, s.password, s.Url)
	if err != nil {
		return err
	}
	s.sftpClient = sftpClient
	ok, err = s.pathExists()
	if err != nil {
		return err
	}
	if !ok {
		err = s.sftpClient.MkdirAll(s.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sftp) Read()  ([]byte, error) {
	fp, err := s.sftpClient.Open("1.html")
	if err != nil {
		fmt.Printf("open remote file :%s, err:%+v failed.\n", "1.html", err)
		return nil, err
	}
	defer fp.Close()
	bytes, err := ioutil.ReadAll(fp)
	return bytes, err
}

func (s *Sftp) Write() error {
	for k, item := range items {
		fPath := path.Join(s.Path, k)
		writer, err := s.sftpClient.OpenFile(fPath, os.O_RDWR|os.O_CREATE)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(item))
		if err != nil {
			return err
		}
		writer.Close()
	}

	return nil
}


func (s *Sftp) Close() error {
	err := s.sftpClient.Close()
	if err != nil {
		return err
	}
	return nil
}

func main()  {
	f := NewSftpProvider()
	_ = f.Connect()

	B,_:=f.Read()
	fmt.Println(string(B))
	_ = f.Write()
	_ = f.Close()
}