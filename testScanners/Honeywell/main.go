package main

import (
	"HelloWorld/testScanners/Honeywell/Com"
	"HelloWorld/testScanners/Honeywell/GoUsb"
	"HelloWorld/testScanners/Honeywell/StandardKey"
	"fmt"
	"runtime"
)

type Scanner interface {
	Connect()
	Read() (string, error)
}

var scanner Scanner
var Env string

const (
	StandardInput= "as a standard input"
	ComIn  = "com in"
	GousbIn  = "gousb in"
)
func init() {
	Env=GousbIn
}

func main()  {
	switch Env {
	case StandardInput:
		scanner=&StandardKey.Standard{}
	case ComIn:
		if runtime.GOOS == "windows" {
			scanner=&Com.Com{Name: "COM3"}
			break
		}
		scanner=&Com.Com{Name: "/dev/ttyACM0"}
	case GousbIn:
		if runtime.GOOS == "windows" {
			break
		}
		scanner=&GoUsb.GoUsb{Name: "3118:2305"}
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
