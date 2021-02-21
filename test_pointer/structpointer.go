package main

import (
	"fmt"
	"time"
)

type User struct {
	name string
	age  int
}

//ch1里面的值传递 不因改变而改变
//ch2 引用传递 改变即改变
//ch3  传结构体则值传递，传指针则引用传递
func main() {
	ch1 := make(chan User, 1000)        //创建一个无缓存channel
	ch2 := make(chan *User, 1000)       //创建一个无缓存channel
	ch3 := make(chan interface{}, 1000) //创建一个无缓存channel
	m := User{
		name: "杨敬伟",
		age:  18,
	}
	t := &m
	ch1 <- *t
	ch1 <- m //往通道写数据

	ch2 <- &m
	ch3 <- &m //往通道写数据
	ch3 <- m  //往通道写数据

	m.age = 31

	ch1 <- m
	ch2 <- &m
	ch3 <- &m
	ch3 <- m //往通道写数据
	go func() {
		for num := range ch1 {
			fmt.Println(num, "ch1")
		}
	}()

	go func() {
		for num := range ch2 {
			fmt.Println(num, "ch2")
		}
	}()

	go func() {
		for num := range ch3 {
			switch n := num.(type) {
			case User:
				fmt.Println(n)

			case *User:
				fmt.Println(n)

			}
			fmt.Println(num, "ch3")
		}
	}()

	time.Sleep(2 * time.Second)
	close(ch1)
	close(ch2)
	close(ch3)
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
