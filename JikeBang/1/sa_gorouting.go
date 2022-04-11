package main

import "fmt"

//https://studygolang.com/articles/22390?fr=sidebar
//import (
//	"fmt"
//	"os"
//	"time"
//)
//
//func loop() {
//	for i := 0; i <10 ; i++ {
//		fmt.Printf("%d + %v \n", i, os.Getpid())
//		time.Sleep(1*time.Second)
//	}
//}
//
//func main() {
//	go loop() // 启动一个goroutine
//	fmt.Println("hello")
//	loop()
//	fmt.Println("world")
//	for {
//		time.Sleep(1*time.Second)
//	}
//}

//
//func main() {
//
//	messages := make(chan string)
//
//	go func() { messages <- "ping"
//	fmt.Println("hello")
//
//	}()
//	time.Sleep(1*time.Second)
//	fmt.Println(runtime.NumGoroutine())
//	msg := <-messages
//
//	fmt.Println(msg)
//	for  {
//		time.Sleep(1*time.Second)
//	}
//}

//func main() {
//fmt.Println(runtime.NumCPU())
//	fmt.Println(runtime.GOMAXPROCS(8))
//	fmt.Println(runtime.NumGoroutine())
//}

//var c = make(chan int, 10)
//var a string
//
//func f() {
//	a = "hello, world"
//close(c)
//}
//
//func main() {
//	go f()
//	fmt.Println(<-c)
//	print(a)
//}

var a = []int{1, 2, 3}
var b = map[int]string{1: "aa",
	2: "bb"}

var c = make(map[int]string)
var d = make([]string, 10)

var e = []int{4, 5, 6}

func main() {
	var c = make(chan int)
	go func() {
		c <- 1
		c <- 2
		c <- 3
		close(c)
	}()
	//for v := range c {
	//	fmt.Println(v)
	//}

	for {

		fmt.Println(<-c)
	}
}
