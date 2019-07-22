/*
func main() {
	resp, err := http.Post("http://www.baidu.com",
		"application/x-www-form-urlencoded",
		strings.NewReader("username=xxx&password=xxxx"))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
*/
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

func main() {

	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		// 错误处理
		log.Println(err)
		return
	}

	defer resp.Body.Close() //关闭链接

	headers := resp.Header

	for k, v := range headers {
		fmt.Printf("k=%v, v=%v\n", k, v) //所有头信息
	}

	fmt.Printf("resp status %s,statusCode %d\n", resp.Status, resp.StatusCode)

	fmt.Printf("resp Proto %s\n", resp.Proto)

	fmt.Printf("resp content length %d\n", resp.ContentLength)

	fmt.Printf("resp transfer encoding %v\n", resp.TransferEncoding)

	fmt.Printf("resp Uncompressed %t\n", resp.Uncompressed)

	fmt.Println(reflect.TypeOf(resp.Body))
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	length, _ := buf.ReadFrom(resp.Body)
	fmt.Println(len(buf.Bytes()))
	fmt.Println(length)
	fmt.Println(string(buf.Bytes()))
}
