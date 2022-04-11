package io

import (
	"fmt"
	"sync"
	"time"

	"github.com/masami10/rush/services/device"
	"github.com/pkg/errors"
)

type IOModule struct {
	device.BaseDevice
	config *ConfigIO
	client IO

	flashInterval time.Duration
	closing       chan struct{}
	flashes       map[uint16]uint16
	mtx           sync.Mutex
	opened        bool
	diag          Diagnostic
}

func NewIOModule(flashInterval time.Duration, conf *ConfigIO, d Diagnostic, service *Service, serialNumber string) *IOModule {
	s := &IOModule{
		diag:          d,
		config:        conf,
		flashInterval: flashInterval,
		opened:        false,
		flashes:       map[uint16]uint16{},
		closing:       make(chan struct{}, 1),
	}
	s.BaseDevice = device.CreateBaseDevice(device.BaseDeviceTypeIO, d, service, serialNumber)
	return s
}

func (s *IOModule) SetIONotify(notify IONotify) {}

func (s *IOModule) SetSerialNumber(sn string) {
	s.BaseDevice.SetSerialNumber(sn)
	s.config.SN = sn
}

func (s *IOModule) getIONotifyService() IONotify {
	return s.BaseDevice.GetParentService().(IONotify)
}

func (s *IOModule) WillStart() error {
	if vendor, ok := VendorModels[s.config.Model]; !ok {
		return errors.Errorf("Model: %s Is Not Support", s.config.Model)
	} else {
		switch vendor.Type() {
		case IoModbustcp:
			s.client = &ModbusTcp{
				cfg:    s.config,
				notify: s.getIONotifyService(),
				vendor: vendor,
			}
		default:
			return errors.New(fmt.Sprintf("invalid model type: %s", s.config.Model))
		}
	}

	return s.BaseDevice.WillStart()
}

func (s *IOModule) Start() error {
	if err := s.BaseDevice.Start(); err != nil {
		return err
	}

	if err := s.WillStart(); err != nil {
		return err
	}

	go s.flashProc()
	s.opened = true

	return s.client.Start()
}

func (s *IOModule) Model() interface{} {
	return nil
}

func (s *IOModule) Stop() error {
	if s.opened {
		s.closing <- struct{}{}

		return s.client.Stop()
	}
	return nil
}

func (s *IOModule) Status() string {
	return s.client.Status()
}

func (s *IOModule) IORead() (string, string, error) {
	return s.client.IORead()
}

func (s *IOModule) IOWrite(index uint16, status uint16) error {
	switch status {
	case OutputStatusOff:
		// 从flash列表删除
		s.removeFlash(index)
	case OutputStatusFlash:
		// 加入flash列表
		s.addFlash(index)
		return nil
	}

	if s.client == nil {
		return errors.New("IO Client Is Nil")
	}

	return s.client.IOWrite(index, status)
}

func (s *IOModule) DeviceType() string {
	return "io"
}

func (s *IOModule) Children() map[string]device.IBaseDevice {
	return map[string]device.IBaseDevice{}
}

func (s *IOModule) Data() interface{} {
	inputs, outputs, err := s.IORead()
	if err != nil {
		return nil
	}

	return IoData{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (s *IOModule) Config() interface{} {
	vendor := VendorModels[s.config.Model]

	return IoConfig{
		InputNum:  vendor.Cfg().InputNum,
		OutputNum: vendor.Cfg().OutputNum,
	}
}

func (s *IOModule) addFlash(output uint16) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.flashes[output] = output
}

func (s *IOModule) removeFlash(output uint16) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.flashes, output)
}

func (s *IOModule) getFlashes() map[uint16]uint16 {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.flashes
}

func (s *IOModule) flashProc() {
	flashTicker := time.NewTicker(s.flashInterval)
	defer flashTicker.Stop()
	status := OutputStatusOff
	flag := -1
	for {
		select {
		case <-flashTicker.C:
			// 状态0<->1变换
			flag *= -1
			status += flag

			flashes := s.getFlashes()
			for _, v := range flashes {
				err := s.IOWrite(v, uint16(status))
				if err != nil {
					s.diag.Error("IOWrite Failed", err)
				}
			}

		case <-s.closing:
			return
		}
	}
}
