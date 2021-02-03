package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

type UserInfo struct {
	Name   string
	Gender string
	Age    int
}

func kua(arg string) (string, error) {
	return arg + "真帅", nil
}
func main() {
	// 解析指定文件生成模板对象
	htmlByte, err := ioutil.ReadFile("./testTemplate/func_template/hello.html")
	if err != nil {
		fmt.Println("read html failed, err:", err)
		return
	}
	// 自定义一个夸人的模板函数
	/*
		kua := func(arg string) (string, error) {
			return arg + "真帅", nil
		}
	*/
	// 采用链式操作在Parse之前调用Funcs添加自定义的kua函数
	tmpl, err := template.New("hello").Funcs(template.FuncMap{"kua": kua}).Parse(string(htmlByte))
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}

	user := UserInfo{
		Name:   "小明",
		Gender: "男",
		Age:    18,
	}
	// 使用user渲染模板，并将结果写入w
	tmpl.Execute(os.Stdout, user)
}
