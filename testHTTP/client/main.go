package main

import (
	"fmt"
	"net/http"
)

//File:HTTP persistent connection.svg
//
//http 1.0中默认是关闭的，需要在http头加入"Connection: Keep-Alive"，
// 才能启用Keep-Alive；http 1.1中默认启用Keep-Alive，如果加入"Connection: close "，
// 才关闭。目前大部分浏览器都是用http1.1协议，也就是说默认都会发起Keep-Alive的连接请求了，
// 所以是否能完成一个完整的Keep-Alive连接就看服务器设置情况。
func main() {
	//resp, err := http.Get("http://www.baidu.com")
	resp, err := http.Get("http://127.0.0.1:12345/home")
	if err != nil {
		fmt.Println("http.Get err = ", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Status = ", resp.Status)
	fmt.Println("StatusCode = ", resp.StatusCode)
	fmt.Println("Header = ", resp.Header)
	//fmt.Println("Body = ", resp.Body)

	buf := make([]byte, 4*1024)
	var tmp string

	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("read err = ", err)
			break
		}
		tmp += string(buf[:n])
	}

	//读取网页内容，打印出来
	fmt.Println("tmp = ", tmp)

	//buf:=make([]byte,1)

	//n,_:=os.Stdin.Read(buf)hhhh

	//fmt.Println(string(buf[:n]))
	///h,_:=fmt.Scan(buf)
	//	fmt.Println(h)
	//fmt.Println(string(buf))
	///http.Post()
	///	println(os.Stdin)
	//resp1,_:=http.Get("http://www.baidu.com")
	//var a []string
	//	io.Copy(os.Stdout,resp1.Body)

	fmt.Println("dfdfgthghghlkjjjjjjjjjjjjjjjjjjjjjjjjjj")
	//body, _ := ioutil.ReadAll(resp1.Body)

	///fmt.Println(string(body))
	//var name  string
	//fmt.Scanf("%s",&name)
	//fmt.Printf("%s",name)

}
