package main

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"time"
)

func main() {

	dateString := "2019-10-16T11:20:30+08:00"
	dateString2 := "2019-10-16 11:20:30"


	dt := time.Now()
	fmt.Println(dt)

	fmt.Println(dt.UTC())
	fmt.Println(dt.Local())
	fmt.Println(dt.Location())

	if dateString2 != "" {
		dt11, _ := time.Parse("2006-01-02 15:04:05", dateString2)
		fmt.Println("dt11",dt11)
	}

	if dateString != "" {
		dt1, _ := time.Parse(time.RFC3339, dateString)
		fmt.Println(dt1)
	}


	if dateString != "" {
		loc, _ := time.LoadLocation("Local")
		dt, _ = time.ParseInLocation(time.RFC3339, dateString, loc)
		dt12, _ := time.ParseInLocation("2006-01-02 15:04:05", dateString2, loc)
		fmt.Println("dt12",dt12)
	}

	//UpdateTime := dt.UTC()
	fmt.Println(dt)
	fmt.Println(dt.Format(time.RFC3339))


	fmt.Println("special time",time.Now().Format("06010215 04"))

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
