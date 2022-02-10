package main

import (
	"fmt"
	"math/rand"
	"time"
)

var a = "123456"

func main() {
	n := 4
	buf := [100]byte{58, 54, 49, 52, 53, 54}
	fmt.Println(string(buf[:n]))
	rand.Seed(time.Now().UnixNano())

	for {
		for i := 0; i < n; i++ {
			buf[i] = byte(rand.Intn(128))
		}
		fmt.Println(string(buf[:6]))
		if a == string(buf[:6]) {
			fmt.Println(fmt.Sprintf("hello%s", string(buf[:6])))
			return
		}

	}
}
