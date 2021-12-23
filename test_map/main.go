package main

import (
	"fmt"
	"log"
)

func main() {
	m := make(map[int32]string)
	fmt.Println(len(m),cap([]int64{}))
	m[0] = "EDDYCJY1"
	fmt.Println(len(m))
	m[1] = "EDDYCJY2"
	m[2] = "EDDYCJY3"
	m[3] = "EDDYCJY4"
	m[4] = "EDDYCJY5"
	fmt.Println(len(m))
	for k, v := range m {
		log.Printf("k: %v, v: %v", k, v)
	}
}
