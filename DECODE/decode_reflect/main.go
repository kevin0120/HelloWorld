package main

import (
	"HelloWorld/DECODE/ascii"
	"fmt"
	"strings"
)

const (
	TEST_STRINGS = "003" + "3.14" + "7410" + "15.3" + "-897" + "99" + "1" + "0" + "sn001" + "666"
)

type Header struct {
	TOOL string `start:"1" end:"5"`
	Sn   int    `start:"6" end:"8"`
}

type OpenProtocolHeader struct {
	L   int64   `start:"1" end:"3"`
	LEN float32 `start:"4" end:"7"`
	MID string  `start:"8" end:"11"`
	MD  float64 `start:"12" end:"15"`
	M   int     `start:"16" end:"19"`
	Faa uint    `start:"20" end:"21"`
	B   bool    `start:"22" end:"22"`
	C   bool    `start:"23" end:"23"`
	TO  Header  `start:"24" end:"..."`
}

func main() {
	//var he =Header{}
	fmt.Println(fmt.Sprintf("Test Data: %# 20X", TEST_STRINGS))
	var testop = OpenProtocolHeader{}
	err := ascii.Unmarshal(TEST_STRINGS, &testop)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", testop)

	//一个容量为o的切片并不一定==nil！！！！！
	/*
	   var	test map[int][]*OpenProtocolHeader
	   	test=map[int][]*OpenProtocolHeader{}

	   	if _, ok := test[23]; !ok {
	   		test[23] = nil
	   	}

	   	fmt.Println("hhhhhhhhhhh",len(test[23]))
	   	test[23]= append(test[23], &testop)

	   	test[23]=test[23][1:]

	   	if test[23]!=nil{
	   		fmt.Println("hhhhhhhhhhh",len(test[23]))
	   	}

	*/

	fmt.Println(strings.TrimSpace("3. 4"))
}
