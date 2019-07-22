package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	/*互斥锁1 互斥锁
	传统并发程序对共享资源进行访问控制的主要手段，由标准库代码包中sync中的Mutex结构体表示。*/
	var mutex sync.Mutex
	fmt.Println("start lock main")
	mutex.Lock()
	fmt.Println("get locked main")
	for i := 1; i <= 3; i++ {
		go func(i int) {
			fmt.Println("start lock ", i)
			/*如果对一个已经上锁的对象再次上锁，那么就会导致该锁定操作被阻塞，直到该互斥锁回到被解锁状态*/
			mutex.Lock()
			fmt.Println("get locked ", i)
			mutex.Unlock()
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Println("Unlock the lock main")
	mutex.Unlock()
	fmt.Println("get unlocked main")
	time.Sleep(time.Second)
	fmt.Println("#####################################")
	/*2 读写锁
	读写锁是针对读写操作的互斥锁，可以分别针对读操作与写操作进行锁定和解锁操作 。
	读写锁的访问控制规则如下：
	①多个写操作之间是互斥的
	②写操作与读操作之间也是互斥的
	③多个读操作之间不是互斥的
	在这样的控制规则下，读写锁可以大大降低性能损耗。
	由标准库代码包中sync中的RWMutex结构体表示
	*/
	var rwm sync.RWMutex
	for i := 0; i < 3; i++ {
		go func(i int) {
			fmt.Println("try to lock read ", i)
			rwm.RLock()
			fmt.Println("get locked ", i)
			time.Sleep(time.Second * 2)
			fmt.Println("try to unlock for reading ", i)
			rwm.RUnlock()
			fmt.Println("unlocked for reading ", i)
		}(i)
	}
	time.Sleep(time.Millisecond * 1000)
	fmt.Println("try to lock for writing")
	rwm.Lock()
	fmt.Println("locked for writing")

}
