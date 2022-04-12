package main

import (
	"fmt"
	"os"
	"path"
)

const (
	OpenFilePerm = os.FileMode(0644)
)

var items = map[string]string{"1": "hello ftp", "2": "kevin"}

type File struct {
	Path string
}

func NewFileProvider() *File {
	return &File{
		Path: "./testFTP/rush/localfile",
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func openFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, OpenFilePerm)
}

func (f *File) Connect() error {
	ok, err := PathExists(f.Path)
	if err != nil {
		return err
	}
	if !ok {
		err = os.MkdirAll(f.Path, OpenFilePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *File) Close() error {
	return nil
}

func (f *File) Write() error {
	for k, item := range items {
		fPath := path.Join(f.Path, k)
		writer, err := openFile(fPath)
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
func (f *File) Read() ([]byte, error) {
	var message []byte
	var m []byte
	for k, _ := range items {
		fPath := path.Join(f.Path, k)
		reader, err := openFile(fPath)
		if err != nil {
			return message, err
		}

		n, err := reader.Read(m)
		message = append(message, m[:n]...)

		if err != nil {
			return message, err
		}

		reader.Close()
	}
	return message, nil
}

func main() {
	f := NewFileProvider()
	_ = f.Connect()

	B, _ := f.Read()
	fmt.Println(string(B))
	_ = f.Write()
	_ = f.Close()
}
