package main

import "fmt"

func main()  {
	var param string = "日"

	// string to []byte
	b := []byte(param)
	fmt.Println(b)

	// []byte to string
	s := string(b)
	fmt.Println(s)
}
