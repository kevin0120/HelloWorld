package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	serverAddr := "192.168.60.40:8000"

	// 172.28.0.180
	//localIP := []byte{0xAC, 0x1C, 0, 0xB4}  // 指定IP
	localIP := []byte{} //  any IP，不指定IP
	localPort := 9001   // 指定端口

	netAddr := &net.TCPAddr{Port: localPort}

	if len(localIP) != 0 {
		netAddr.IP = localIP
	}

	fmt.Println("netAddr:", netAddr)

	d := net.Dialer{Timeout: 5 * time.Second, LocalAddr: netAddr}
	conn, err := d.Dial("tcp", serverAddr)

	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buffer := make([]byte, 512)
	reader := bufio.NewReader(conn)

	n, err2 := reader.Read(buffer)
	if err2 != nil {
		fmt.Println("Read failed:", err2)
		return
	}

	fmt.Println("count:", n, "msg:", string(buffer))

	select {}
}
