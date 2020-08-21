package main

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"time"
)

func main() {

	dateString := "2019-10-16T11:20:30+08:00"

	dt := time.Now()
	fmt.Println(dt)

	fmt.Println(dt.UTC())
	fmt.Println(dt.Local())
	fmt.Println(dt.Location())

	if dateString != "" {
		loc, _ := time.LoadLocation("Local")
		dt, _ = time.ParseInLocation(time.RFC3339, dateString, loc)
	}

	//UpdateTime := dt.UTC()
	fmt.Println(dt)
	fmt.Println(dt.Format(time.RFC3339))
	//指定时间
	fmt.Println(time.Unix(11, 01).UTC())

	//时间戳
	fmt.Println(time.Unix(11, 01).Unix())
	fmt.Println(time.Unix(66666611, 01).UnixNano())

	errors.Prefix="kevin"
	e1:=errors.New("hello world")

	fmt.Println(e1)
	fmt.Println(e1.Error())




}
