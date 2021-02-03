package main

import (
	"fmt"
	"time"

	"context"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	ch := make(chan interface{})

	go func() {

		time.Sleep(10 * time.Second)
		ch <- 1
	}()
	close(ch)

	select {
	case <-ch:
		fmt.Println("ch")
		return

	case <-ctx.Done():
		fmt.Println("done")
		return

	}

}

//select 没有default时一直阻塞直到有chanel中有数据取出
