package main

import (
	"HelloWorld/testScanners/Honeywell/GoUsb"
	"fmt"
)


func main()  {
	//var scanner=&StandardKey.Stand{
	//}

	//var scanner=&Com.Com{Name: "COM3"}
	//if runtime.GOOS != "windows" {
	//	//demesg find com
	//	//sudo usermod -aG　dialout kevin  -----get powe
	//	scanner=&Com.Com{Name: "/dev/ttyACM0"}
	//}

	var scanner=&GoUsb.GoUsb{Name: "3118:2305"}

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
