// Copyright 2018-2020 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gopcua/opcua/uacp"
	"log"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "opc.tcp://localhost:53530/OPCUA/SimulationServer", "OPC UA Endpoint URL")
	)
	flag.Parse()

	ctx := context.Background()

	log.Printf("Listening on %s", *endpoint)
	l, err := uacp.Listen(*endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	c, err := l.Accept(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	log.Printf("conn %d: connection from %s", c.ID(), c.RemoteAddr())

	go func() {

		for {
			b, _ := c.Receive()

			fmt.Println(b)
		}

	}()

	select {}

}
