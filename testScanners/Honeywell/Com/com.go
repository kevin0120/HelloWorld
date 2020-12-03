package Com

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"time"
)

type Com struct {
	Name string
	port    *serial.Port
}

func (s1 *Com) Read() (string, error) {
	//// 写入字符串“012345”
	//n, err := s.Write([]byte("012345"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////延时100
	//time.Sleep(100)

	buf := make([]byte, 128)
	n, err := s1.port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Read %d Bytes\r\n", n)
	//for i := 0; i < n; i++ {
	//	fmt.Printf("buf[%d]=%c\r\n", i, buf[i])
	//}
	return string(buf[0:n]), nil
}

func (s1 *Com) Connect() {

	//设置串口编号
	c := &serial.Config{Name: s1.Name, Baud: 115200}

	//打开串口
	for {
		s, err := serial.OpenPort(c)
		if err != nil {
			log.Fatal(err)
			time.Sleep(1 * time.Second)
		}
		s1.port = s
		break
	}

	fmt.Printf("%s 连接成功\n", s1.Name)
}
