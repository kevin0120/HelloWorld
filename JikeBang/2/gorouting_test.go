package main


import (
	"fmt"
	"testing"
)

func Test_1(t *testing.T) {
	c := make(chan struct{})
	go func() {
		sum := 0
		for i := 0; i < 100; i++ {
			sum += i
		}
		fmt.Println("jieguoshi", sum)
		c <- struct{}{}
		fmt.Println("hello")
		fmt.Println("1")
		fmt.Println("2")
		fmt.Println("3")
		fmt.Println("1")
		fmt.Println("2")
		fmt.Println("3")
		fmt.Println("1")
		fmt.Println("2")
		fmt.Println("3")
		for  {
			fmt.Println("yjw")
			//time.Sleep(1*time.Second)
		}
	}()
	<-c
	fmt.Println("world")
}
