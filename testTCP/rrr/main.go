package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
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
	service := "127.0.0.1:8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	fmt.Println(tcpAddr)

	// fmt.Println("远程地址:")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	//fmt.Println(conn)err
	for {
		var msr [512]byte
		/*
			_, err = conn.Write([]byte("GGGGGGGGGGGGGG"))
			checkError(err)
		*/
		//result, err := ioutil.ReadAll(conn)
		//fmt.Println(string(msr[0]))
		//os.Exit(0)
		_, err = conn.Write([]byte("hello"))
		if err!=nil{
			fmt.Println(err)
		}
		_, err = conn.Read(msr[0:])
		if err!=nil{
			fmt.Println(err)
		}
		fmt.Println("recieve :", string(msr[0:]))
		time.Sleep(10 * time.Second)
	}
	//err := conn.Close()

}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "cuowu:%s", err.Error())
		os.Exit(1)
	}
}
