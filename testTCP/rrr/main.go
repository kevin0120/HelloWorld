package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	fmt.Println(time.Now())
	/*
		if len(os.Args)!=2{
			fmt.Fprintf(os.Stderr,"用法:%s ip地址\n",os.Args[0])
			os.Exit(1)
		}
		fmt.Println(len(os.Args))
		name:=os.Args[1]
		addr:=net.ParseIP(name)
		if addr==nil{
			fmt.Println("wuxiao")
		}else {
			fmt.Println("IP:",addr.String())
		}
	*/
	//if len(os.Args) != 2 {
	//	fmt.Fprintf(os.Stderr, "用法:%s ip地址\n", os.Args[0])
	//	os.Exit(1)
	//}
	service := "10.1.1.54:8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	fmt.Println(tcpAddr)

	// fmt.Println("远程地址:")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
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
	go func() {
		var msr [512]byte
		for  {
				n, err := conn.Read(msr[0:])
				if err != nil {
					fmt.Printf("%v read error :%s \n", time.Now(), err)
					conn.Close()
					break
				} else {
					fmt.Printf("%v recieve :%s \n", time.Now(), string(msr[0:n]))
				}
		}
	}()

	//time.Sleep(5 * time.Second)
	//conn.Close()


i:=0
	for {
		//var msr [512]byte
		//////
		/////*
		////	_, err = conn.Write([]byte("GGGGGGGGGGGGGG"))
		////	checkError(err)
		////*/
		//////result, err := ioutil.ReadAll(conn)
		//////fmt.Println(string(msr[0]))
		//////os.Exit(0)
		////go func() {
		//	_, err = conn.Read(msr[0:])
		//	if err != nil {
		//		fmt.Printf("%v read error :%s \n", time.Now(), err)
		//	} else {
		//		fmt.Printf("%v recieve :%s \n", time.Now(), string(msr[0:]))
		//	}
		//}()
		//
		//
		//
		//_, err = conn.Write([]byte("hello"))
		_, err = conn.Write([]byte(fmt.Sprintf("hello%5d",i)))
i++
		if err != nil {
			fmt.Printf("%v write error :%s \n", time.Now(), err)
			conn.Close()
			break
		}
		//conn.Close()
		//if err != nil {
		//	fmt.Printf("%v write11111 error :%s \n", time.Now(), err)
		//	if err, ok := err.(net.Error); !ok || !err.Temporary() {
		//		// permanent error. close the connection
		//		conn.Close()
		//		//conn= nil
		//		fmt.Printf("%v write22222 error :%s \n", time.Now(), err)
		//	}
		//}
		time.Sleep(3000 * time.Millisecond)
		//_, err = conn.Write(nil)

	}
	//err := conn.Close()

}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "cuowu:%s", err.Error())
		os.Exit(1)
	}
}
