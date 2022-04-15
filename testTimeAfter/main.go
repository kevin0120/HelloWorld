package main

import (
	"fmt"
	"time"
)

func main() {
	var my_ticker = time.NewTicker(5 * time.Second)
	t1 := time.NewTimer(time.Second * 3)
	for {
		select {
		//只有在本次select 操作中会有效， 再次select 又会重新开始计时,如果两个time after 那么时间大的可能不会执行
		case <-time.After(2 * time.Second):
			fmt.Println("111111111111111111111111111")
		case <-time.After(3 * time.Second):
			fmt.Println("-------------------------------")
		case <-my_ticker.C:
			fmt.Println("222222222222222222222222222")
		case <-t1.C:
			fmt.Println("333333333333333333333333333")

		//default:
		//	fmt.Println("ddddddddddddddddddddddddddd")
		}



		//time.Sleep(1*time.Second)
	}

}
