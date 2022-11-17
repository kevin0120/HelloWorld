package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
	//"time"
)

func echo(conn *net.TCPConn) {
	//for{
	//fmt.Println("远程地址:", conn.RemoteAddr())
	//
	//time.Sleep(3*time.Second)
	//}
	//tick := time.Tick(5 * time.Second)
	for {
		//n,err:=conn.Write([]byte(now.String()))
		//if err!=nil{
		//	log.Println(err)
		//	conn.Close()
		//	return
		//}
		var msr [512]byte
		_, err := conn.Read(msr[0:])
		if err != nil {
			fmt.Println(err)
			//err = conn.Close()

			//break
			return
		}

		fmt.Println("recieve:", string(msr[0:]))
		n, err := conn.Write([]byte("world"))
		fmt.Printf("%v send %d bytes to %s\n", time.Now(), n, conn.RemoteAddr())
		if err != nil {
			fmt.Println(err)
			err = conn.Close()

			//break
			return
		}
	}
}

//"192.168.4.188"
func main() {
	address := net.TCPAddr{
		IP: net.ParseIP(":"),
		//等价与0.0.0.0
		Port: 8011,
	}
	listener, err := net.ListenTCP("tcp4", &address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}

		//ya ce  quanlianjie  duilie
		time.Sleep(3 * time.Millisecond)
		fmt.Println("远程地址:", conn.RemoteAddr())
		//go echo(conn)
		go testConnectServer(conn)
	}

}

func testConnectServer(conn *net.TCPConn) {
	//err := conn.SetKeepAlive(false)
	//if err != nil {
	//	return
	//}
	//err := conn.SetKeepAlivePeriod(10 * time.Second)
	//if err != nil {
	//	return
	//}
	sockFile, sockErr := conn.File()
	if sockErr == nil {
		// got socket file handle. Getting descriptor.
		fd := int(sockFile.Fd())
		// Ping amount
		err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPCNT, 3)
		if err != nil {
			fmt.Println("on setting keepalive probe count", err.Error())
		}
		// Retry interval
		err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_KEEPINTVL, 5)
		if err != nil {
			fmt.Println("on setting keepalive retry interval", err.Error())
		}
		// don't forget to close the file. No worries, it will *not* cause the connection to close.
		sockFile.Close()
	}

	for {
		var msr [512]byte
		n, err := conn.Read(msr[0:])
		if err != nil {
			fmt.Printf("%v read error :%s \n", time.Now(), err)
			err = conn.Close()
			return
		} else {
			fmt.Printf("%v read  :%s \n", time.Now(), msr[0:n])
			conn.Write([]byte("hello world!!"))
			conn.Close()
			return
		}
	}
}
