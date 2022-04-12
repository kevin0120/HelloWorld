package io

import (
	"errors"
	"fmt"
	"time"

	"github.com/simonvetter/modbus"
	"go.uber.org/atomic"
)

const (
	TIMEOUT        = 3 * time.Second
	DefaultReadItv = 300 * time.Millisecond
)

type ModbusTcp struct {
	cfg    *ConfigIO
	status atomic.Value

	inputs  atomic.Value
	outputs atomic.Value

	client *modbus.ModbusClient

	vendor Vendor
	notify IONotify
}

func (s *ModbusTcp) Start() error {
	s.status.Store(IoStatusOffline)
	s.inputs.Store("")
	s.outputs.Store("")

	go s.connect()

	return nil
}

func (s *ModbusTcp) Stop() error {
	if s.client != nil {
		return s.client.Close()
	}

	return nil
}

func (s *ModbusTcp) Status() string {
	return s.status.Load().(string)
}

func (s *ModbusTcp) createClient() (err error) {
	if s.client != nil {
		_ = s.client.Close()
	}

	if c, e := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     s.cfg.Address,
		Speed:   s.cfg.Speed,
		Timeout: TIMEOUT,
	}); e != nil {
		err = e
		return
	} else {
		s.client = c
	}
	return
}

func (s *ModbusTcp) connect() {
	if err := s.createClient(); err != nil {
		panic(fmt.Sprintf("modbustcp connect  error: %s", err.Error()))
	}

	for {
		err := s.client.Open()

		if err != nil {
			continue
		} else {
			// online
			s.status.Store(IoStatusOnline)
			s.notify.OnStatus(s.cfg.SN, IoStatusOnline)
			s.read()
		}
	}
}

func (s *ModbusTcp) read() {
	readItv := time.Duration(s.cfg.ReadItv)
	if readItv == 0 {
		readItv = DefaultReadItv
	}

	ticker := time.NewTicker(readItv)
	defer ticker.Stop()

	for {
		<-ticker.C
		_, _, err := s.IORead()
		if err == nil {
			continue
		}
		// offline
		s.status.Store(IoStatusOffline)
		s.notify.OnStatus(s.cfg.SN, IoStatusOffline)
		go s.connect() //重连
		return

	}
}

func (s *ModbusTcp) formatIO(results []bool) string {
	rt := ""
	for i := 0; i < len(results); i++ {
		if results[i] {
			rt += "1"
			continue
		}
		rt += "0"
	}

	return rt
}

func (s *ModbusTcp) IORead() (string, string, error) {
	client := s.client
	var err error
	var result []bool

	inputs := ""
	outputs := ""

	if client == nil {
		return inputs, outputs, errors.New("client is nil")
	}

	// input status
	if s.vendor.Cfg().InputNum > 0 {
		switch s.vendor.Cfg().InputReadType {
		case ReadTypeDiscretes:
			result, err = client.ReadDiscreteInputs(s.vendor.Cfg().InputAddress, s.vendor.Cfg().InputNum)
		}

		if err != nil {
			return inputs, outputs, err
		}

		inputs = s.formatIO(result)
	}

	// output status
	if s.vendor.Cfg().OutputNum > 0 {
		switch s.vendor.Cfg().OutputReadType {
		case ReadTypeCoils:
			result, err = client.ReadCoils(s.vendor.Cfg().OutputAddress, s.vendor.Cfg().OutputNum)
		}

		if err != nil {
			return inputs, outputs, err
		}

		outputs = s.formatIO(result)
	}

	if s.inputs.Load().(string) != inputs {
		s.inputs.Store(inputs)
		s.notify.OnChangeIOStatus(s.cfg.SN, IoTypeInput, inputs)
	}

	if s.outputs.Load().(string) != outputs {
		s.outputs.Store(outputs)
		s.notify.OnChangeIOStatus(s.cfg.SN, IoTypeOutput, outputs)
	}

	return inputs, outputs, nil
}

func (s *ModbusTcp) IOWrite(index uint16, status uint16) error {
	var bstatus bool
	if s.Status() == IoStatusOffline {
		return errors.New(IoStatusOffline)
	}

	if index > (s.vendor.Cfg().OutputNum - 1) {
		return errors.New("invalid index")
	}

	var err error
	client := s.client

	if status == 0 {
		bstatus = false
	} else {
		bstatus = true
	}

	switch s.vendor.Cfg().WriteType {
	case WriteTypeSingleCoil:
		err = client.WriteCoil(index, bstatus)
	}

	return err
}

func (s *ModbusTcp) SetIONotify(notify IONotify) {
	s.notify = notify
}
