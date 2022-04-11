package main

import (
	"fmt"
	"time"
)

func main()  {
	var my_ticker = time.NewTicker(5 * time.Second)


	for  {
		select {
		case <- time.After(2*time.Second):
			fmt.Println("111111111111111111111111111")

		case <- my_ticker.C:
			fmt.Println("222222222222222222222222222")

		}
		//time.Sleep(1*time.Second)
	}


}
