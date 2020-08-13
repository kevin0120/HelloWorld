package main

import (
	"fmt"
	"time"
)

//传入channel后不变了
func main() {
	ch := make(chan *int, 1000) //创建一个无缓存channel

	var a = 20 /* 声明实际变量 */
	var b = 30
	var ip *int /* 声明指针变量 */

	ip = &a  /* 指针变量的存储地址 */
	ch <- &a //往通道写数据

	ch <- ip

	ip = &b
	ch <- ip

	go func() {
		for num := range ch {
			fmt.Println("num = ", *num)
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
