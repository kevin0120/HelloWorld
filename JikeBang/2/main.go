package main

import "fmt"

func main()  {
	c := make(chan struct{})
	go func() {
		sum := 0
		for i := 0; i < 100; i++ {
			sum += i
		}
		fmt.Println("jieguoshi", sum)
		c <- struct{}{}
		fmt.Println("hello")
	}()
	<-c
	fmt.Println("world")
}