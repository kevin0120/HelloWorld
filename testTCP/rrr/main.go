package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	localAddr := net.TCPAddr{
		IP: net.ParseIP(":"),
		//等价与0.0.0.0
		Port: 8888,
	}

	serverAddr := net.TCPAddr{
		IP: net.ParseIP("192.168.60.40"),
		//等价与0.0.0.0
		Port: 8000,
	}

	fmt.Println("netAddr:", localAddr, serverAddr)

	service := "192.168.60.40:8000"

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	fmt.Println(tcpAddr)

	// fmt.Println("远程地址:")

	//localAddr为nil的话则随机本地ip+端口
	conn, err := net.DialTCP("tcp", &localAddr, &serverAddr)
	checkError(err)
	_, err = conn.Write([]byte("GGGGGGGGGGGGGG"))
	testConnectClient(conn)

	//conn.SetReadBuffer(0)
	//fmt.Println(conn)err
	//go func() {
	//	select {
	//	case <- time.After(15*time.Second):
	//		conn.Close()
	//	}
	//
	//
	//
	//go func() {
	//	var msr [512]byte
	//	for {
	//		n, err := conn.Read(msr[0:])
	//		if err != nil {
	//			fmt.Printf("%v read error :%s \n", time.Now(), err)
	//			conn.Close()
	//			break
	//		} else {
	//			fmt.Printf("%v recieve :%s \n", time.Now(), string(msr[0:n]))
	//		}
	//	}
	//}()
	//
	//time.Sleep(70 * time.Second)
	//conn.Close()

	//i := 0
	//for {
	//	//var msr [512]byte
	//	//////
	//	/////*
	//	////	_, err = conn.Write([]byte("GGGGGGGGGGGGGG"))
	//	////	checkError(err)
	//	////*/
	//	//////result, err := ioutil.ReadAll(conn)
	//	//////fmt.Println(string(msr[0]))
	//	//////os.Exit(0)
	//	////go func() {
	//	//	_, err = conn.Read(msr[0:])
	//	//	if err != nil {
	//	//		fmt.Printf("%v read error :%s \n", time.Now(), err)
	//	//	} else {
	//	//		fmt.Printf("%v recieve :%s \n", time.Now(), string(msr[0:]))
	//	//	}
	//	//}()
	//	//
	//	//
	//	//
	//	//_, err = conn.Write([]byte("hello"))
	//	_, err = conn.Write([]byte(fmt.Sprintf("hello%5d", i)))
	//	i++
	//	if err != nil {
	//		fmt.Printf("%v write error :%s \n", time.Now(), err)
	//		conn.Close()
	//		break
	//	}
	//	//conn.Close()
	//	//if err != nil {
	//	//	fmt.Printf("%v write11111 error :%s \n", time.Now(), err)
	//	//	if err, ok := err.(net.Error); !ok || !err.Temporary() {
	//	//		// permanent error. close the connection
	//	//		conn.Close()
	//	//		//conn= nil
	//	//		fmt.Printf("%v write22222 error :%s \n", time.Now(), err)
	//	//	}
	//	//}
	//	time.Sleep(3000 * time.Millisecond)
	//	//_, err = conn.Write(nil)
	//
	//}
	////err := conn.Close()

}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "cuowu:%s", err.Error())
		os.Exit(1)
	}
}

func testConnectClient(conn *net.TCPConn) {
	//err1 := conn.SetKeepAlive(true)
	//if err1 != nil {
	//	return
	//}
	//err2 := conn.SetKeepAlivePeriod(5 * time.Second)
	//if err2 != nil {
	//	return
	//}
	go func() {
		for {
			fmt.Println("111111111111")
			var msr [512]byte
			_, err := conn.Read(msr[0:])
			if err != nil {
				fmt.Printf("%v read error :%s \n", time.Now(), err)
				conn.Close()

				break

			} else {
				fmt.Printf("%v recieve :%s \n", time.Now(), string(msr[0:]))
			}
		}
	}()
	//go func() {
	//	for {
	//		fmt.Println("222222222222")
	//		time.Sleep(30 * time.Second)
	//		_, err := conn.Write([]byte("fffffffffffffffff"))
	//		if err != nil {
	//			fmt.Printf("%v write error :%s \n", time.Now(), err)
	//			conn.Close()
	//
	//			break
	//
	//		}
	//	}
	//}()
	time.Sleep(70000 * time.Second)
	err := conn.Close()
	if err != nil {
		return
	}
}
