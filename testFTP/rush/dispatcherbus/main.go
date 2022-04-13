package main

import (
	"fmt"
	"github.com/kevin0120/HelloWorld/testFTP/rush/dispatcherbus/dispatcherbus"
	"github.com/kevin0120/HelloWorld/testFTP/rush/dispatcherbus/utils"
	"time"
)

func Func1(s interface{}) {
	fmt.Printf("This is Func1,with data %s\n", s)
}
func Func2(s interface{}) {
	fmt.Printf("This is Func2,with data %s\n", s)
}
func Func3(s interface{}) {
	fmt.Printf("This is Func3,with data %s\n", s)
}
func Func4(s interface{}) {
	fmt.Printf("This is Func4,with data %s\n", s)
}
func Func5(s interface{}) {
	fmt.Printf("This is Func5,with data %s\n", s)
}

var diagnostic dispatcherbus.Diagnostic

//1. create, register tart dispatch
//
//2.Create一定要在最前面,并且只允许运行一次即可，多次后面的跳过
//3.register可以多次 程序轮询运行 ,很多在hmi服务里
//4.dispatch 用来触发
//5.

func main() {
	dispatcherBus, _ := dispatcherbus.NewService(diagnostic)

	dispatcherMap := dispatcherbus.DispatcherMap{
		dispatcherbus.DispatcherServiceStatus: utils.CreateDispatchHandlerStruct(Func1),
		dispatcherbus.DispatcherResult:        utils.CreateDispatchHandlerStruct(Func1),
	}

	// create, register and start
	dispatcherBus.LaunchDispatchersByHandlerMap(dispatcherMap)


	_ = dispatcherBus.Create(dispatcherbus.DispatcherDeviceStatus, utils.DefaultDispatcherBufLen)
	_ = dispatcherBus.Create(dispatcherbus.DispatcherReaderData, utils.DefaultDispatcherBufLen)

	//Register 和 Start 没有先后顺序
	// 接收设备状态变化
	dispatcherBus.Register(dispatcherbus.DispatcherDeviceStatus, utils.CreateDispatchHandlerStruct(Func2))
	_ = dispatcherBus.Start(dispatcherbus.DispatcherDeviceStatus)
	// 接收读卡器数据
	_ = dispatcherBus.Start(dispatcherbus.DispatcherReaderData)
	dispatcherBus.Register(dispatcherbus.DispatcherReaderData, utils.CreateDispatchHandlerStruct(Func2))
	dispatcherBus.Register(dispatcherbus.DispatcherReaderData, utils.CreateDispatchHandlerStruct(Func2))
	dispatcherBus.Register(dispatcherbus.DispatcherReaderData, utils.CreateDispatchHandlerStruct(Func2))


	_ = dispatcherBus.Dispatch(dispatcherbus.DispatcherReaderData, "hello"+dispatcherbus.DispatcherServiceStatus)

	for {
		time.Sleep(1 * time.Second)
	}
	//dispatcherBus.ReleaseDispatchersByHandlerMap(dispatcherMap)
}
