package main

import (
	"HelloWorld/go_system/drivers"
	"fmt"
	"os"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	//if runtime.GOOS == "windows" {
	//	// 依赖的dll libs
	//	libFiles := []string{"hello.txt"}
	//
	//	box := packr.NewBox("./dll")
	//
	//	for _, l := range libFiles {
	//		if PathExists(l) {
	//			continue
	//		}
	//		data, err := box.Find(l)
	//		if err == nil {
	//			ioutil.WriteFile(l, data, fs.ModePerm)
	//		}
	//	}
	//	pathOld := os.Getenv("PATH")
	//	pwd, _ := os.Getwd()
	//	if !strings.Contains(pathOld, pwd) {
	//		fmt.Println("hh")
	//		//os.Setenv("PATH", fmt.Sprintf("%s;%s", pathOld, pwd))
	//	}
	//}
	drivers.Nkio()
	fmt.Println("hello world!!")
}
