package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

/*Done 方法在Context被取消或超时时返回一个close的channel,close的channel可以作为广播通知，告诉给context相关的函数要停止当前工作然后返回。
当一个父operation启动一个goroutine用于子operation，这些子operation不能够取消父operation。下面描述的WithCancel函数提供一种方式可以取消新创建的Context.
Context可以安全的被多个goroutine使用。开发者可以把一个Context传递给任意多个goroutine然后cancel这个context的时候就能够通知到所有的goroutine。
Err方法返回context为什么被取消。
Deadline返回context何时会超时。
Value返回context相关的数据。
*/
/*BackGound是所有Context的root，不能够被cancel。*/
/*WithCancel返回一个继承的Context,这个Context在父Context的Done被关闭时关闭自己的Done通道，或者在自己被Cancel的时候关闭自己的Done。
WithCancel同时还返回一个取消函数cancel，这个cancel用于取消当前的Context。

作者：kingeasternsun
链接：https://www.jianshu.com/p/d24bf8b6c869
来源：简书
简书著作权归作者所有，任何形式的转载都请联系作者获得授权并注明出处。*/

var logg *log.Logger

func someHandler() {
	ctx, cancel := context.WithCancel(context.Background())
	go doStuff(ctx)

	//10秒后取消doStuff
	time.Sleep(10 * time.Second)
	cancel()

}

//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
func doStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

/**************************************************************/
func timeoutHandler() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	// go doTimeOutStuff(ctx)
	go doStuff(ctx)

	time.Sleep(10 * time.Second)

	cancel()

}

func timeoutHandler1() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	go doTimeOutStuff(ctx)
	//go doStuff(ctx)

	time.Sleep(10 * time.Second)

	cancel()

}

func doTimeOutStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)

		if deadline, ok := ctx.Deadline(); ok { //设置了deadl
			logg.Printf("deadline set")
			if time.Now().After(deadline) {
				logg.Printf(ctx.Err().Error())
				return
			}

		}

		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

func main() {
	fmt.Println(fmt.Sprintf("downdddddddddddddddddddddddd"))
	logg = log.New(os.Stdout, "", log.Ltime)
	someHandler()
	logg.Printf("down")
	//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
	fmt.Println("************************************************")
	timeoutHandler()
	//每1秒work一下，同时会判断ctx是否超时，如果是就退出
	fmt.Println("************************************************")
	//可以采集deadline的时间
	timeoutHandler1()

}
