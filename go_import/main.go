package main

import (
	_ "HelloWorld/go_import/first"
	_ "HelloWorld/go_import/third"
	"fmt"
	"time"
)

var _ = func() error {
	fmt.Printf("main.go下第一个全局变量运行时间%s\n", time.Now())
	return nil
}()

func init() {
	fmt.Printf("main.go下init时间%s\n", time.Now())
}

var _ = func() error {
	fmt.Printf("main.go下第二个全局变量运行时间%s\n", time.Now())
	return nil
}()

func main() {

	fmt.Printf("main.go下main函数时间%s\n", time.Now())

	fmt.Println("hello world")
}
