package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int) //创建一个无缓存channel

	//新建一个goroutine
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i //往通道写数据
		}
		//不需要再写数据时，关闭channel
		//close(ch)
		//ch <- 666 //关闭channel后无法再发送数据

		time.Sleep(1 * time.Second)
		ch <- 111
		time.Sleep(1 * time.Second)
		close(ch)
	}() //别忘了()
	time.Sleep(10 * time.Second)
	for num := range ch {
		fmt.Println("num = ", num)
	}
	//读取已经关闭的通道并不会阻塞,而是返回类型零值.关闭时ok==false
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
