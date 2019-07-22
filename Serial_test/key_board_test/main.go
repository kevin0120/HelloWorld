package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	file2, _ := os.Open("/dev/input/event21")
	var a []byte
	go func() {
		for {
			file2.ReadAt(a, 5)
			fmt.Println(a)
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(30 * time.Second)

}
