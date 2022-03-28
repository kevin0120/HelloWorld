package test

// Copyright 2018 Trung Hieu Le. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.
import (
	"log"
	"os"
	"testing"
	"time"

	"fmt"
	"github.com/robinson/gos7"
	"sync"
)

const (
	tcpDevice = "192.168.5.99"
	rack      = 0
	slot      = 1
)
// slot 2 for 300/400, slot 1 for 1200/1500
// conn.initiateConnection({port: 102, host: '192.168.0.2', localTSAP: 0x0100, remoteTSAP: 0x0200, timeout: 8000, doNotOptimize: true}, connected);
// local and remote TSAP can also be directly specified instead. The timeout option specifies the TCP timeout.

func TestTCPClient(t *testing.T) {
	handler := gos7.NewTCPClientHandler(tcpDevice, rack, slot)
	handler.Timeout = 200 * time.Second
	handler.IdleTimeout = 200 * time.Second
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	handler.Connect()
	defer handler.Close()
	client := gos7.NewClient(handler)
	ClientTestAll(t, client)
}

func TestMultiTCPClient(t *testing.T) {
	var handlers sync.Map
	var clients sync.Map

	tcpDevices := make([]map[string]string, 2)
	tcpDevices[0] = make(map[string]string, 1)
	tcpDevices[1] = make(map[string]string, 1)
	tcpDevices[0]["tcpDevice"] = "192.168.10.19:102"
	tcpDevices[1]["tcpDevice"] = "192.168.10.10:102"

	c := make(chan int)

	for k := range tcpDevices {
		go func(device map[string]string) {
			handler := gos7.NewTCPClientHandler(tcpDevice, rack, slot)
			handler.Timeout = 200 * time.Second
			handler.IdleTimeout = 200 * time.Second
			handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
			handler.Address = device["tcpDevice"]
			handler.Connect()
			handlers.Store(device["tcpDevice"], handler)

			client := gos7.NewClient(handler)
			clients.Store(device["tcpDevice"], client)

			c <- 1
		}(tcpDevices[k])
	}

	var cS []int

	for n := range c {
		cS = append(cS, n)
		if len(cS) == len(tcpDevices) {
			close(c)
			break
		}
	}

	cli, exist := clients.Load("192.168.10.10:102")
	client, ok := cli.(gos7.Client)
	if exist && ok {
		buf := make([]byte, 255)
		client.AGReadDB(200, 34, 4,  buf)
		var s7 gos7.Helper
		var result float32
		s7.GetValueAt(buf, 0, &result)
		fmt.Printf("%v\n", result)
	}

	defer func() {
		handlers.Range(func(key, value interface{}) bool {
			h, _ := handlers.Load(key)
			if hh, ok := h.(*gos7.TCPClientHandler); ok {
				hh.Close()
			}
			return true
		})
	}()
}