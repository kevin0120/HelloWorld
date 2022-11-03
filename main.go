package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

func main() {
	a := strings.TrimRight(string("hello "), "\x20")
	fmt.Println(a)

	b := []byte("abcdfscdff")

	fmt.Println(len(b))
	c := bytes.IndexAny(b, "cc")

	fmt.Println(c)

	d := time.Now()
	fmt.Println(d)
	var e interface{}
	e = "2022-07-06 11:14:40.000000 +00:00"
	f := e.(time.Time)

	fmt.Println(f)
}
