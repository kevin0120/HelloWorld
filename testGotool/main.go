package main

import "fmt"

//go tool compile -S ./main.go >> main.S

func main()  {
	i:=0
	i++
	fmt.Println(i)
}
