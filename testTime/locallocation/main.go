package main

import (
	"fmt"
	"syscall"
	"time"
)

func main() {
	env, _ := syscall.Getenv("ZONEINFO")

	fmt.Println(env)

	t := time.Now().UTC()
	fmt.Println(t)
	local3, _ := time.LoadLocation("Asia/Shanghai")

	t1 := t.Local().Format("2006-01-02T15:04:05")

	t2 := t.In(local3).Format("2006-01-02T15:04:05")

	fmt.Println(t1)
	fmt.Println(t2)
}
