package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)


var b = flag.Bool("b", false, "bool类型参数")
var s = flag.String("s", "", "string类型参数")


func main () {
	for idx, args := range os.Args {
		fmt.Println("参数" + strconv.Itoa(idx) + ":", args)
	}


	flag.Parse()
	fmt.Println("-b:", *b)
	fmt.Println("-s:", *s)
	fmt.Println("其他参数：", flag.Args())
}


//go run main.go -b  -s test others
//D:\Code\lianxi\HelloWorld\go_argsFlag>go run main.go -b -s test others
//参数0: C:\Users\admin\AppData\Local\Temp\go-build427800672\b001\exe\main.exe
//参数1: -b
//参数2: -s
//参数3: test
//参数4: others
//-b: true
//-s: test
//其他参数： [others]

