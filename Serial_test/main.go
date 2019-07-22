package main

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {

	//设置串口编号
	c := &serial.Config{Name: "/dev/input/event21", Baud: 115200}

	//打开串口
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// 写入字符串“012345”
	n, err := s.Write([]byte("012345"))
	if err != nil {
		log.Fatal(err)
	}

	//延时100
	time.Sleep(100)

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Read %d Bytes\r\n", n)
	for i := 0; i < n; i++ {
		fmt.Printf("buf[%d]=%c\r\n", i, buf[i])
	}
}
