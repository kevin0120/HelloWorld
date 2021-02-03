package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

//rang 里的i 是共享的如果range下多个协程访问i会错乱,例如print(a)改成print(i)
//如果要吧waitGROUP 当参数传入函数,一定要传指针!!!!

//state1[0]: wait() numbers
//
//state1[1]: Add() Done() operation area
//
//state1[2]:用于同步线程的  阻塞 runtime_Semacquire(semap) 释放 runtime_Semrelease(semap, false, 0)
//
//释放次数是wait（）次数

func main() {

	wg := sync.WaitGroup{}
	//si := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	wg.Add(1)
	//for _,i:= range si {
	//	wg.Add(1)
	//	go func(a int) {
	//		println(a)
	//		wg.Done()
	//	}(i)
	//}

	go func() {
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	wg.Wait()
	//
	//fmt.Println("哈哈哈,我活干完了")
	var state1 [3]uint32
	state1[0] = 1
	state1[1] = 1
	fmt.Println(unsafe.Pointer(&state1[0]))
	fmt.Println(unsafe.Pointer(&state1[1]))
	fmt.Println(*(*uint64)(unsafe.Pointer(&state1[1])))
	fmt.Println(unsafe.Pointer(&state1[0]))

	c := uintptr(unsafe.Pointer(&state1[0]))

	if c%8 == 0 {
		a := (*uint64)(unsafe.Pointer(&state1))
		b := &state1[2]

		fmt.Println(*a, *b)
	} else {
		a := (*uint64)(unsafe.Pointer(&state1[1]))
		b := &state1[0]
		fmt.Println(*a, *b)
	}

}
