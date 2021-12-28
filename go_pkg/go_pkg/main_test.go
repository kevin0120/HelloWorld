package main

import (
	"fmt"
	"math/rand"
)

func main() {
	sum := 0.0
	n := 10000000
	i := 0
	for i < n {
		x := rand.ExpFloat64()
		sum += x
		i += 1
	}
	expect := sum / (float64)(n)
	fmt.Println(expect)
}
