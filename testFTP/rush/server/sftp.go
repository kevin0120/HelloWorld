package file

import (
	"github.com/masami10/rush/services/outputs/serializer"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path"
	"time"
)

type Sftp struct {
	Url  string
	Path string

	username   string
	password   string
	sftpClient *sftp.Client
	serializer serializer.Serializer
}

func NewSftpProvider(config *Plugin) Provider {
	return &Sftp{
		Url:      config.Url,
		Path:     config.Path,
		username: config.User,
		password: config.PassWord,
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

func (s *Sftp) Write(pkg []serializer.PublishPackage) error {
	items, err := s.serializer.Serialize(pkg)
	if err != nil {
		return err
	}
	for _, item := range items {
		fPath := path.Join(s.Path, item.Filename)
		writer, err := s.sftpClient.OpenFile(fPath, os.O_RDWR|os.O_CREATE)
		if err != nil {
			return err
		}
		_, err = writer.Write(item.Data)
		if err != nil {
			return err
		}
		writer.Close()
	}

	return nil
}

func (s *Sftp) SetSerializer(serializer serializer.Serializer) {
	s.serializer = serializer
}

func (s *Sftp) Close() error {
	err := s.sftpClient.Close()
	if err != nil {
		return err
	}
	return nil
}
