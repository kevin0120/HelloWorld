package driver

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"io/fs"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func init() {
	if runtime.GOOS == "windows" {
		// 依赖的dll libs
		libFiles := []string{"libspc_lib.dll"}

		box := packr.New("mybox", "../../cmake-build-debug")

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
			os.Setenv("PATH", fmt.Sprintf("%s;%s", pathOld, pwd))

			pathOld = os.Getenv("PATH")
			//fmt.Println(pathOld)
		}
	}
}

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
