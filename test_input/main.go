package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {

	fmt.Println("请输入一个字符串：")
	//读键盘
	reader := bufio.NewReader(os.Stdin)
	//以换行符结束
	str, _ := reader.ReadString('\n')
	fmt.Println(str)

	log.Print(0xc2e, 0x0901)
	for {
		fmt.Println("***********")
		var a string
		fmt.Scanln(&a)
		fmt.Println(a)
	}
}
