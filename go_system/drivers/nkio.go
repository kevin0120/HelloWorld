package drivers

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"io/fs"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
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

func Nkio() {
	if runtime.GOOS == "windows" {
		//_ = packr.PackJSONBytes("../cmake-build-debug", "1.txt", "\"MTIzU1NTU1M=\"")
		// 依赖的dll libs
		libFiles := []string{"1.txt"}

		box := packr.New("hello", "../dll")

		for _, l := range libFiles {
			if PathExists(l) {
				continue
			}
			data, err := box.Find(l)
			if err == nil {
				_ = ioutil.WriteFile(l, data, fs.ModePerm)
			}
		}

		pathOld := os.Getenv("PATH")
		pwd, _ := os.Getwd()
		if !strings.Contains(pathOld, pwd) {
			fmt.Println("hh")
			//os.Setenv("PATH", fmt.Sprintf("%s;%s", pathOld, pwd))
		}
	}
}
