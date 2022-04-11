package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

func main() {
	//操作系统
	fmt.Println("GOOS:", runtime.GOOS)
	//架构
	fmt.Println("GOARCH:", runtime.GOARCH)
	//GOROOT
	fmt.Println("GOROOT:", runtime.GOROOT())
	//go版本
	fmt.Println("Version:", runtime.Version())
	//cpu数
	fmt.Println("NumCPU:", runtime.NumCPU())
	fmt.Println()

	//MAC和IP地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	for _, inter := range interfaces {
		fmt.Println(inter.Name, inter.HardwareAddr)

		addrs, _ := inter.Addrs()
		for _, addr := range addrs {
			fmt.Println("  ", addr.String())
		}
	}

	name, _ := os.Hostname()
	fmt.Println(name)
}
