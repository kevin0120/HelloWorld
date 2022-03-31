package main

import (
	"fmt"
	"github.com/masami10/spc/v2/wrapper/golang/spc"
)

func main() {
	var a = []float64{5., 5., 10., 12., 5., 5., 10., 12., 5., 5.}

	c, _ := spc.Cpk(a, 10, 14, 4)
	fmt.Printf("cpk: %f", c)
}
