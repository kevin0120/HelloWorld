package main

import (
	"fmt"
	"time"
)

//传入channel后改变map channel 内的值一起改变
func main() {
	ch := make(chan map[string]string, 1000) //创建一个无缓存channel

	m := make(map[string]string)
	m["china"] = "beijing"

	ch <- m //往通道写数据

	m["janpan"] = "dingjing"
	ch <- m
	ch <- m
	go func() {
		for num := range ch {
			fmt.Println(num)
			num["test"] ="hello"

		}

	}()

	time.Sleep(5 * time.Second)
	close(ch)
	time.Sleep(1 * time.Second)
	//读取已经关闭的通道并不会阻塞,而是返回类型零值.关闭时ok==false
	//forange时不会rang到
	/*
		for{
			v,ok :=<-ch
			if ok {
				fmt.Println("##############")
			}else{
				break
			}
		}
	*/

	//time.Sleep(100*time.Second)
}
