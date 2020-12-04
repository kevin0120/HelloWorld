package Com

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"time"
)


//　　dmesg find com,　sudo chmod 777 /dev/ttyACM0
//修改权限为可读可写可执行，但是这种设置电脑重启后，又会出现这种问题，还要重新设置．因此查询资料，可以用下面这条指令：
//
//　　sudo usermod -aG　dialout kevin

//windows下需要在官网上下载驱动



type Com struct {
	Name string
	port    *serial.Port
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
			continue
		}
		s1.port = s
		break
	}

	fmt.Printf("%s 连接成功\n", s1.Name)
}

func (s1 *Com) Read() (string, error) {
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
