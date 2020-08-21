package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("hello world")
	fmt.Println(fmt.Errorf("%s", "hello err").Error())

	//此为在标准输出打印字符串并不能赋值给变量
	n, _ := fmt.Printf("%s+%d", "hello format", 1224)
	fmt.Println(n)
	//下面这个才是拼接字符传的函数
	fmt.Println(fmt.Sprintf("%s+%d", "hello format", 1224))
	//吧其他格式转换成字符串
	fmt.Println(len(fmt.Sprint(111)))
	//将整数转换成asii码
	fmt.Println(string(88))
	//string转成int：
	int, _ := strconv.Atoi("9999")
	//string转成int64 base代表进制，bitsize代表内存大小：
	int64, _ := strconv.ParseInt("8888", 10, 64)
	///int转成string：
	string1 := strconv.Itoa(9999)
	///int64转成string：
	string2 := strconv.FormatInt(7777, 10)
	fmt.Println(int, int64, string1, string2)

}