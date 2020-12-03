package main

import (
	"HelloWorld/testScanners/Honeywell/Com"
	"fmt"
)





func main()  {
	//var scanner=&StandardKey.Stand{
	//}
    var scanner=&Com.Com{Name: "Com3"}
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
