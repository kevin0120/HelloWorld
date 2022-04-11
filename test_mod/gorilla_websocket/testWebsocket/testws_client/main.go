package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/rush/v1/ws"}
	log.Printf("connecting to %s", u.String())
	fmt.Println(u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
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
			hh := strings.Split(st, ":")
			bb := strings.Replace(st, hh[0]+":", "", -1)
			jj := strings.Split(bb, ";")
			//fmt.Println("收到一个message:", jj[2])

			fmt.Println(fmt.Sprintf("收到第%d个message:%s", i, jj[2]))
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

	for i := 0; i < 1; i++ {
		//	//请求获取工单列表
		//	order := `{
		//	"sn": 1000,
		//	"type": "WS_ORDER_LIST",
		//	"data": {
		//	"time_from": "",
		//		"time_to": "",
		//		"status": "",
		//		"page_size": 2,
		//		"page_no": 2
		//}
		//}`
		//	_ = c.WriteMessage(websocket.TextMessage, []byte(order))
		//
		//	//请求获取工单详情
		//	order = `{
		//	"sn": 2000,
		//	"type": "WS_ORDER_DETAIL",
		//	"data": {
		//		"id": 1
		//}
		//}`
		//	_ = c.WriteMessage(websocket.TextMessage, []byte(order))
		//
		//	//根据code和workcenter取得工单
		//	order = `{
		//	"sn":3000,
		//	"type":"WS_ORDER_DETAIL_BY_CODE",
		//	"data":{
		//		"code":"MO00033",
		//		"workcenter":""
		//}
		//}`
		//	_ = c.WriteMessage(websocket.TextMessage, []byte(order))
		//
		//	//更新工单状态
		//	order = `{
		//	"sn":4000,
		//	"type":"WS_ORDER_UPDATE",
		//	"data":{
		//		"id":8,
		//		"status":"delete"
		//}
		//}`
		//	_ = c.WriteMessage(websocket.TextMessage, []byte(order))
		//
		//	//更新工步状态
		//	order = `{
		//	"sn":5000,
		//	"type":"WS_ORDER_STEP_UPDATE",
		//	"data":{
		//		"id":8,
		//		"status":"delete"
		//}
		//}`
		//	_ = c.WriteMessage(websocket.TextMessage, []byte(order))
		//
		//	//更新工步数据
		//	order = `{
		//	"sn":6000,
		//	"type":"WS_ORDER_STEP_DATA_UPDATE",
		//	"data":{
		//		"id":8,
		//		"data":"hhhhh"
		//}
		//}`
		//	_ = c.WriteMessage(websocket.TextMessage, []byte(order))
		//
		//	//开工请求
		//	order = `{
		//	"sn":7000,
		//	"type":"WS_ORDER_START_REQUEST",
		//	"data":{
		//		"code":"MO0001",
		//		"track_code":"XM205Z03",
		//		"product_code":"M000001780589",
		//		"workcenter": "TA2-26L-01",
		//		"date_start": "2019-10-16T11:20:30+08:00",
		//		"resources": {
		//			"users": [
		//						{
		//							"name": "张三",
		//							"code": "111-333-5555"
		//						}
		//					],
		//			"equipments": [
		//							{
		//									"name": "数显扳手#1",
		//									"code": "33322-111"
		//							}
		//							]
		//					}
		//			}
		//}`
		//	_ = c.WriteMessage(websocket.TextMessage, []byte(order))

		//拧紧工具使能
		order := `{"sn":57487841693,"data":{"controller_sn":"0002","enable":true,"tool_sn":"xx0011"},"type":"WS_TOOL_ENABLE"}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(order))

		order = `{"sn":57487841693,"data":{"controller_sn":"c1","enable":true,"tool_sn":"xx4443"},"type":"WS_TOOL_ENABLE"}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(order))

		//拧紧工具Pset
		order = `{"sn":57487577813,"data":{"controller_sn":"0002","pset":1,"sequence":1,"step_id":4,"tool_sn":"xx0011","total":3,"user_id":"1","workorder_id":2},"type":"WS_TOOL_PSET"}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(order))

		order = `{"sn":57487577813,"data":{"controller_sn":"c1","pset":2,"sequence":1,"step_id":4,"tool_sn":"xx4443","total":3,"user_id":"1","workorder_id":2},"type":"WS_TOOL_PSET"}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(order))

		order = `{"sn":69197138450,"type":"WS_DEVICE_STATUS"}`
		_ = c.WriteMessage(websocket.TextMessage, []byte(order))

	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write-read:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write-read close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
