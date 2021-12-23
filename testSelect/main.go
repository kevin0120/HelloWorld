package main

import (
	"fmt"
	"time"
)
//发送者
func sender(c chan int) {
	for i := 0; i < 100; i++ {
		c <- i
		if i >= 5 {
			time.Sleep(time.Second * 1)
		} else {
			time.Sleep(time.Second)
		}
	}
}


func main() {
	c := make(chan int)
	go sender(c)
	//timeout := time.After(time.Second * 3)
	for {
		select {
		case d := <-c:
			fmt.Println(d)
		//case <-timeout:
		//	fmt.Println("这是定时操作任务 >>>>>")
		case dd := <-time.After(time.Second * 3):
			fmt.Println(dd, "这是超时*****")
		}

		fmt.Println("for end")
	}
}
