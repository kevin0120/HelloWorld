package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/resty.v1"
	"log"
	"time"
)

var Headers = map[string]string{"1": "ck_god", "2": "god_girl"}

var Body = []byte("hello,I am a http client!")

const (
	URL = "http://127.0.0.1:8082/rush/v1/hello/3?q=1"

	WEBSOCKETURL = "ws://127.0.0.1:8082/rush/v1/ws"
)

func main() {
	go httpclient()
	go websocketclient()

	for {
		time.Sleep(1 * time.Second)
	}
}

func httpclient() {
	var httpClient *resty.Client
	client := resty.New()
	client.SetRESTMode() // restful mode is default
	client.SetTimeout(time.Duration(6 * time.Second))
	client.SetContentLength(true)
	// Headers for all request
	client.SetHeaders(Headers)
	client.
		SetRetryCount(3).
		SetRetryWaitTime(time.Duration(5 * time.Second)).
		SetRetryMaxWaitTime(20 * time.Second)

	httpClient = client

	resp, err := httpClient.R().
		SetBody(Body).
		Get(URL)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}

func websocketclient() {
	c, _, err := websocket.DefaultDialer.Dial(WEBSOCKETURL, nil)

	//	fmt.Println(a)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		i := 1
		for {
			_, message, err := c.ReadMessage()
			st := string(message)
			//hh := strings.Split(st, ":")
			//bb := strings.Replace(st, hh[0]+":", "", -1)
			//jj := strings.Split(bb, ";")
			//fmt.Println("收到一个message:", jj[2])

			fmt.Println(fmt.Sprintf("收到第%d个message:%s", i, st))
			i++
			if err != nil {
				log.Println("read:", err)
				return
			}
			//log.Printf("recv: %s", message)
		}
	}()

	reg := `{
		"sn": 999,
		"type": "WS_REG",
		"data": {
		"hmi_sn": "hhh"
	}
	}`
	_ = c.WriteMessage(websocket.TextMessage, []byte(reg))

	for {
		time.Sleep(time.Second)
	}
}
