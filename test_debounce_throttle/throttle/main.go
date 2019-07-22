package main

import (
	"fmt"
	"github.com/bep/throttle"
	"sync/atomic"
	"time"
)

func main() {
	var counter uint64
	//var a, b int
	f := func() {
		atomic.AddUint64(&counter, 1)
		//	fmt.Println(a, b)
	}

	throttled := throttle.New(10 * time.Millisecond)

	for i := 0; i < 3; i++ {
		for j := 0; j < 100; j++ {
			throttled(f)
			//	a, b = i, j
			time.Sleep(5 * time.Millisecond)
		}

		time.Sleep(20 * time.Millisecond)
	}

	c := int(atomic.LoadUint64(&counter))

	fmt.Println("Counter is", c)
	// Output: Counter is 3
}
