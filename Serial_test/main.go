package main

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"time"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

func main() {

	//设置串口编号
	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 115200}

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
	go func() {
		file, _ := os.OpenFile("./myFile", os.O_WRONLY|os.O_CREATE, 0666)
		for {
			file.WriteString("Hello world!!!!\r\n")
			file.Write([]byte("ts="))
			file.WriteString(string(time.Now().Format(RFC3339Milli)))

			time.Sleep(1 * time.Second)
			log.Println("这是一次测试!!!!", time.Now().String())
		}
		defer file.Close()
	}()
	//延时100
	time.Sleep(100)
	for {
		buf := make([]byte, 128)
		n, err = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Read %d Bytes\r\n", n)
		fmt.Println(buf)
		fmt.Println(string(buf[0:n]))
		/*
			for i := 0; i < n; i++ {
				fmt.Printf("buf[%d]=%c\r\n", i, buf[i])
			}
		*/
	}
}
