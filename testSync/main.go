package main

import (
	"fmt"
	"sync"
	"time"
)

//var once1 sync.Once
func aa() {
	var once1 sync.Once
	once1.Do(func() {
		fmt.Println(time.Now(), "这是第几次调用11：")
	})
}
func main() {
	aa()
	aa()

	var once sync.Once
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(f)
		}()
	}
	time.Sleep(5 * time.Second)
}

func f() {
	fmt.Println(time.Now(), "这是第几次调用：")
}
