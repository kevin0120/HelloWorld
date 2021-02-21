package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestReader(t *testing.T) {

	fmt.Println("请输入一个字符串：")
	//读键盘
	reader := bufio.NewReader(os.Stdin)
	//以换行符结束
	str, _ := reader.ReadString('\n')
	fmt.Println(str)

	log.Print("dfdfdfd")

}
