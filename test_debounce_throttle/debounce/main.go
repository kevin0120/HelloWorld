package main

import (
	"fmt"
	"github.com/bep/debounce"
	"sync/atomic"
	"time"
)

func main() {
	var counter uint64

	f := func() {
		atomic.AddUint64(&counter, 1)
	}

	debounced := debounce.New(100 * time.Millisecond)

	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			debounced(f)
		}

		time.Sleep(200 * time.Millisecond)
	}

	c := int(atomic.LoadUint64(&counter))

	fmt.Println("Counter is", c)
	// Output: Counter is 3
}
