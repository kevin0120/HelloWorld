package io

import (
	"sync"
	"time"
)

type Module struct {
	config        *ConfigIO
	client        *ModbusTcp
	flashInterval time.Duration
	closing       chan struct{}
	flashes       map[uint16]uint16
	mtx           sync.Mutex
	opened        bool
}

func NewIOModule(flashInterval time.Duration, conf *ConfigIO, serialNumber string) *Module {
	s := &Module{
		config:        conf,
		flashInterval: flashInterval,
		opened:        false,
		flashes:       map[uint16]uint16{},
		closing:       make(chan struct{}, 1),
	}
	return s
}

func (s *Module) Model() interface{} {
	return nil
}

func (s *Module) IOWrite(index uint16, status uint16) error {
	switch status {
	case OutputStatusOff:
		// 从flash列表删除
		s.removeFlash(index)
	case OutputStatusFlash:
		// 加入flash列表
		s.addFlash(index)
		return nil
	}
	return nil
}

func (s *Module) DeviceType() string {
	return "io"
}

func (s *Module) Config() interface{} {
	vendor := VendorModels[s.config.Model]

	return IoConfig{
		InputNum:  vendor.Cfg().InputNum,
		OutputNum: vendor.Cfg().OutputNum,
	}
}

func (s *Module) addFlash(output uint16) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.flashes[output] = output
}

func (s *Module) removeFlash(output uint16) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.flashes, output)
}

func (s *Module) getFlashes() map[uint16]uint16 {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.flashes
}

func (s *Module) flashProc() {
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
				}
			}

		case <-s.closing:
			return
		}
	}
}
