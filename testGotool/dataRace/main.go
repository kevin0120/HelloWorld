package main

import (
	"fmt"
	"time"
)

var race int

//go build -race ./main.go


func main() {

	go func() {
		for i := 0; i < 10; i++ {
			race++
			fmt.Printf("gorouting1:%d\n",race)
		}
	}()

	for i := 0; i < 10; i++ {
		race++
		fmt.Printf("gorouting2:%d\n",race)
	}

	time.Sleep(2*time.Second)
	fmt.Println(race)

}
