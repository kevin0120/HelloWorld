package main

import (
	"fmt"
	"log"
	"net"
	//"time"
)

func echo(conn *net.TCPConn) {
	//tick := time.Tick(5 * time.Second)
	for {
		//n,err:=conn.Write([]byte(now.String()))
		//if err!=nil{
		//	log.Println(err)
		//	conn.Close()
		//	return
		//}
		var msr [512]byte
		msr[0] = 'h'
		_, err := conn.Read(msr[0:])
		fmt.Println("远程地址:", string(msr[0]))
		n, err := conn.Write(msr[0:])
		fmt.Printf("send %d bytes to %s\n", n, conn.RemoteAddr())
		if err != nil {
			err = conn.Close()
			fmt.Println(err)
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
		Port: 8000,
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

		fmt.Println("远程地址:", conn.RemoteAddr())
		go echo(conn)

	}

}
