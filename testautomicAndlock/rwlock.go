package main

import (
	"fmt"
	"sync"

	"time"
)

var m *sync.RWMutex

/*读写锁是表示，读的过程中可以有多个读，一旦碰到写的过程只能有一个写，即写的时候只能有一个写并且其他的不能读。*/
func main() {
	fmt.Println(m)
	m = new(sync.RWMutex)

	//写的时候啥都不能干
	go write(1)

	go read(2)

	go write(3)

	time.Sleep(4 * time.Second)

}
func read(i int) {
	fmt.Println(i, "read start")
	m.RLock()
	fmt.Println(i, "reading")
	time.Sleep(1 * time.Second)
	m.RUnlock()
	fmt.Println(i, "read end")
}

func write(i int) {
	fmt.Println(i, "write start")
	m.Lock()
	fmt.Println(i, "writing")
	time.Sleep(1 * time.Second)
	m.Unlock()
	fmt.Println(i, "write end")
}
