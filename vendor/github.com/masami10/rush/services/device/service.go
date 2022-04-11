package device

import (
	"sync"

	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
	"go.uber.org/atomic"

	"github.com/masami10/rush/services/wsnotify"
)

const (
	TIGHTENING_DEVICE_TYPE_TOOL = "tool"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	diag           Diagnostic
	configValue    atomic.Value
	runningDevices map[string]IBaseDevice
	mtxDevices     sync.Mutex
	dispatcherBus  Dispatcher
	dispatcherMap  dispatcherbus.DispatcherMap

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(c Config, d Diagnostic, dp Dispatcher) (*Service, error) {

	s := &Service{
		diag:           d,
		runningDevices: map[string]IBaseDevice{},
		mtxDevices:     sync.Mutex{},
		dispatcherBus:  dp,
	}

	s.configValue.Store(c)

	s.setupGlobalDispatchers()
	s.setupWSRequestHandlers()

	return s, nil
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.initDispatcherRegisters()
	s.dispatcherBus.LaunchDispatchersByHandlerMap(s.dispatcherMap)

	return nil
}

func (s *Service) Close() error {
	s.dispatcherBus.ReleaseDispatchersByHandlerMap(s.dispatcherMap)
	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) initDispatcherRegisters() {
	// 注册websocket请求
	s.dispatcherBus.Register(dispatcherbus.DispatcherWsNotify, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))
}

func (s *Service) setupGlobalDispatchers() {
	s.dispatcherMap = dispatcherbus.DispatcherMap{
		dispatcherbus.DispatcherDeviceStatus:       utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DispatcherReaderData:         utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DispatcherScannerData:        utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DispatcherAnyDeviceInputData: utils.CreateDispatchHandlerStruct(nil),
		dispatcherbus.DispatcherIO:                 utils.CreateDispatchHandlerStruct(nil),
	}
}

func (s *Service) setupWSRequestHandlers() {
	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_DEVICE_STATUS: s.OnWSDeviceStatus,
	})
}

func (s *Service) AddDevice(sn string, d IBaseDevice) {
	defer s.mtxDevices.Unlock()
	s.mtxDevices.Lock()

	_, exist := s.runningDevices[sn]
	if exist {
		return
	}

	s.runningDevices[sn] = d
}

func (s *Service) FetchAllDeviceStatus() []Status {
	return s.fetchAllDevices()
}

func (s *Service) fetchAllDevices() []Status {
	defer s.mtxDevices.Unlock()
	s.mtxDevices.Lock()

	var devices []Status
	for k, v := range s.runningDevices {
		devices = append(devices, Status{
			SN:       k,
			Type:     v.DeviceType(),
			Status:   v.Status(),
			Children: v.Children(),
			Config:   v.Config(),
			Data:     v.Data(),
		})

		for cSN, c := range v.Children() {
			status := Status{
				SN:     cSN,
				Type:   c.DeviceType(),
				Status: c.Status(),
				Config: c.Config(),
				Data:   c.Data(),
			}
			if c.DeviceType() == TIGHTENING_DEVICE_TYPE_TOOL {
				// todo use config sn as tightening unit
				status.TighteningUnit = cSN
			}
			devices = append(devices, status)
		}
	}

	return devices
}
