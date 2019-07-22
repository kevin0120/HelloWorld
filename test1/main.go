package main

import (
	"fmt"
	"os"
)

func main() {
	port0 := os.Getenv("/dev/ttyUSB0")
	port1 := os.Getenv("PORT1")

	fmt.Println(port0, port1)
}
