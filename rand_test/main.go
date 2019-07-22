package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		for {
			a := rand.Intn(100)

			if a > 10 {
				fmt.Println("rand=", time.Now().Unix())
				break
			}
		}
	}

}
