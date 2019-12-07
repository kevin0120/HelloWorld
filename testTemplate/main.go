package main

import (
	"fmt"
	"html/template"
	"os"
)

type UserInfo struct {
	Name   string
	Gender string
	Age    int
}

func main() {
	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("/home/kevin/Downloads/gopath/src/HelloWorld/testTemplate/hello.html")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	// 利用给定数据渲染模板，并将结果写入w
	user := UserInfo{
		Name:   "小明",
		Gender: "男",
		Age:    18,
	}
	tmpl.Execute(os.Stdout, user)
}
