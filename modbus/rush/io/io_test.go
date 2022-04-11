package io

import (
	"fmt"
	"github.com/simonvetter/modbus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIO(t *testing.T) {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://192.168.5.241:502",
		Timeout: 3 * time.Second,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	defer client.Close()

	client.WriteCoil(0, false)
	//client.write
	for {
		//outputs, err := client.ReadCoils(0, 8)
		//if err == nil {
		//	fmt.Println(outputs)
		//} else {
		//	break
		//}

		inputs, err := client.ReadDiscreteInputs(0, 1)
		if err == nil {
			fmt.Print("输出:")
			fmt.Println(inputs)
		} else {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	//_, err = client.WriteMultipleCoils(0, 10, []byte{4, 3})
}

func getIO() *Module {
	cfg := ConfigIO{
		SN:      "1",
		Model:   "MOXA_E1212",
		Address: "192.168.127.201:502",
	}
	return &Module{
		closing: make(chan struct{}, 1),
		client: &ModbusTcp{
			cfg:    &cfg,
			vendor: VendorModels["MOXA_E1212"],
		},

		flashInterval: 1 * time.Second,
	}
}

func TestStart(t *testing.T) {
	io := getIO()
	err := io.client.Start()
	assert.Nil(t, err)
}

func TestStop(t *testing.T) {
	io := getIO()
	err := io.client.Stop()
	assert.Nil(t, err)

	io.client.Start()
	err = io.client.Stop()
	assert.Nil(t, err)
}

func TestWrite(t *testing.T) {
	io := getIO()
	io.client.Start()
	err := io.IOWrite(0, OutputStatusOff)
	assert.NotNil(t, err)

	err = io.IOWrite(0, OutputStatusOn)
	assert.NotNil(t, err)

	err = io.IOWrite(0, OutputStatusFlash)
	assert.Nil(t, err)
}
