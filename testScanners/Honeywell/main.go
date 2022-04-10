package main

import (
	"fmt"
	"github.com/kevin0120/HelloWorld/testScanners/Honeywell/Com"
	"github.com/kevin0120/HelloWorld/testScanners/Honeywell/GoUsb"
	"github.com/kevin0120/HelloWorld/testScanners/Honeywell/StandardKey"
	"log"
	"runtime"
)

type Scanner interface {
	Connect()
	Read() (string, error)
}

var scanner Scanner
var Env string

const (
	StandardInput = "as a standard input"
	ComIn         = "com in"
	GousbIn       = "gousb in"
)

func init() {
	Env = ComIn
}

func main() {
	switch Env {
	case StandardInput:
		scanner = &StandardKey.Standard{}
	case ComIn:
		if runtime.GOOS == "windows" {
			scanner = &Com.Com{Name: "COM3"}
			break
		}
		scanner = &Com.Com{Name: "/dev/ttyACM0"}
	case GousbIn:
		if runtime.GOOS == "windows" {
			log.Fatal("霍尼韦尔的libusb方式暂时不可用")
			break
		}
		scanner = &GoUsb.GoUsb{Name: "3118:2305"}
	}

	scanner.Connect()

	fmt.Println("请使用扫码枪进行扫码")
	for {
		s, err := scanner.Read()

		if err != nil {
			fmt.Println("扫码失败")
		}
		fmt.Printf("成功扫码：%s\n", s)
	}
}
