package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	var gorouting1 int
	var gorouting2 int
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			mu.Lock()
			time.Sleep(99 * time.Millisecond)
			gorouting1++
			fmt.Printf("gorouting1:%d\n",gorouting1)
			mu.Unlock()
		}
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Millisecond)
		mu.Lock()
		gorouting2++
		fmt.Printf("gorouting2:%d\n",gorouting1)
		mu.Unlock()
		time.Sleep(99 * time.Millisecond)
	}

	fmt.Println(gorouting1)
	fmt.Println(gorouting2)
}
