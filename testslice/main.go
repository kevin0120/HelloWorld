package main

import "fmt"

func main() {
	a := []byte("hello")

	fmt.Println(len(a), " ", cap(a))

	fmt.Println(len(a[1:]), " ", cap(a[1:]))
	fmt.Println(a[5:])

	for i := 0; i < len(a); i++ {
		fmt.Println(string(a[i]))
	}
	fmt.Println(string(a))
}
