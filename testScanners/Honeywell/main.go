package main

import (
	"HelloWorld/testScanners/Honeywell/Com"
	"fmt"
	"runtime"
)


func main()  {
	//var scanner=&StandardKey.Stand{
	//}

	//demesg find com
	//sudo usermod -aG　dialout kevin  -----get power
    var scanner=&Com.Com{Name: "/dev/ttyACM0"}
	if runtime.GOOS == "windows" {
		scanner=&Com.Com{Name: "COM3"}
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
