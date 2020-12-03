package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var paths = []string{`C:\workspace\odoo-enterprise`, `C:\workspace\sa_addons`}

func main() {
	for _, p := range paths {
		removePathPyFile(p)
	}
}

func removePathPyFile(paths string) {
	err := filepath.Walk(paths, func(paths string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(paths)

		fileName := filepath.Base(paths)

		fileSuffix := path.Ext(paths)

		filenameOnly := strings.TrimSuffix(fileName, fileSuffix)

		if fileSuffix == ".py" && !isValueInList(filenameOnly, []string{"__init__", "__manifest__", "__openerp__"}) {
			fmt.Println(paths)
			os.Remove(paths)
		}
		//os.Remove(path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func getFileTags(fullFilename string) {
	//fullFilename := "D:/software/Typora/bin/typora.exe"
	fmt.Println("fullFilename =", fullFilename)
	//获取文件名带后缀
	filenameWithSuffix := path.Base(fullFilename)
	fmt.Println("filenameWithSuffix =", filenameWithSuffix)
	//获取文件后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	fmt.Println("fileSuffix =", fileSuffix)

	//获取文件名
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	fmt.Println("filenameOnly =", filenameOnly)

}
