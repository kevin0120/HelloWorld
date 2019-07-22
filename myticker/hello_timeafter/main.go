package main

import (
	"fmt"
	"time"
)

type f func(a ...interface{}) (eorr error)

var my_ticker = time.NewTicker(2 * time.Second)

func test(a ...interface{}) (eorr error) {

	fmt.Println("现在时间是:", a[0], a[1])
	fmt.Println(time.Now())
	return fmt.Errorf("这是现在时间")
}

func flutter(method f, a ...interface{}) (eorr error) {
	var Par []interface{}
	Par = a[0:]
	select {
	case <-my_ticker.C:
		return method(Par[0:]...)
	case <-time.After(1 * time.Second):
		//fmt.Println("这是一个default")
		//fmt.Println(a[0],a[1])
		fmt.Println("已经超时")
	}
	return fmt.Errorf("这是一个default eorr")

}

func main() {
	for i := 0; i < 10; i++ {
		flutter(test, i, "hhh", "ddd")
		//fmt.Println(eorr)
		//test(1,"hhh")

	}
}
