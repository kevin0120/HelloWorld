package io

import (
	"fmt"
	"time"

	"github.com/masami10/rush/services/device"
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/services/httpd"
	"github.com/masami10/rush/services/wsnotify"
	"github.com/masami10/rush/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

type Service struct {
	configValue   atomic.Value
	ios           map[string]IO
	diag          Diagnostic
	dispatcherBus Dispatcher
	deviceService IDeviceService
	httpd         HTTPService
	IONotify

	// websocket请求处理器
	wsnotify.WSRequestHandlers
}

func NewService(c Config, d Diagnostic, dp Dispatcher, ds IDeviceService, httpd HTTPService) *Service {

	s := &Service{
		diag:          d,
		dispatcherBus: dp,
		deviceService: ds,
		ios:           map[string]IO{},
		httpd:         httpd,
	}

	s.configValue.Store(c)

	s.setupWSRequestHandlers()
	s.setupHttpRoute()
	s.loadModules()

	return s
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) GetIOSerialNumberByIdx(index int) (string, int) {
	ii := 0
	for _, v := range s.config().IOS {
		vendor, exist := VendorModels[v.Model]
		if !exist {
			return "", 0
		}

		outputNum := int(vendor.Cfg().OutputNum)
		for i := 0; i < outputNum; i++ {
			if ii == index {
				return v.SN, i
			}
			ii++
		}
	}

	return "", 0
}

func (s *Service) Open() error {
	if !s.config().Enable {
		return nil
	}

	s.initDispatcherRegisters()
	s.initModules()

	return nil
}

func (s *Service) Close() error {

	for _, dev := range s.ios {
		if err := dev.Stop(); err != nil {
			s.diag.Error("Stop io Module Failed ", err)
		}
	}

	return nil
}

func (s *Service) setupHttpRoute() {
	var r httpd.Route = httpd.Route{
		RouteType:   httpd.ROUTE_TYPE_HTTP,
		Method:      "PUT",
		Pattern:     "/io-set",
		HandlerFunc: s.putIOSet,
	}
	if err := s.httpd.AddNewHttpHandler(r); err != nil {
		s.diag.Error("setupHttpRoute error", err)
	}
}

func (s *Service) initDispatcherRegisters() {
	// 注册websocket请求
	s.dispatcherBus.Register(dispatcherbus.DispatcherWsNotify, utils.CreateDispatchHandlerStruct(s.HandleWSRequest))
}

func (s *Service) setupWSRequestHandlers() {
	s.WSRequestHandlers = wsnotify.WSRequestHandlers{
		Diag: s.diag,
	}

	s.SetupHandlers(wsnotify.WSRequestHandlerMap{
		wsnotify.WS_IO_CONTACT: s.OnWSIOContact,
		wsnotify.WS_IO_STATUS:  s.OnWSIOStatus,
		wsnotify.WS_IO_SET:     s.OnWSIOSet,
	})
}

func (s *Service) AddModule(sn string, io IO) {
	if io == nil {
		return
	}

	io.SetIONotify(s)
	s.ios[sn] = io
	s.deviceService.AddDevice(sn, io.(device.IBaseDevice))
}

func (s *Service) loadModules() {
	if !s.config().Enable {
		return
	}
	cfgs := s.config().IOS
	for i := range cfgs {
		v := cfgs[i]
		io := NewIOModule(time.Duration(s.config().FlashInteval), &v, s.diag, s, v.SN)
		s.AddModule(v.SN, io)
	}
}

func (s *Service) initModules() {
	for i := range s.ios {
		io := s.ios[i]
		err := io.Start()
		if err != nil {
			s.diag.Error("start io failed", err)
		}
	}
}

func (s *Service) Read(sn string) (string, string, error) {
	m, err := s.getIO(sn)
	if err != nil {
		return "", "", err
	}

	return m.IORead()
}

func (s *Service) Write(sn string, index uint16, status uint16) error {
	m, err := s.getIO(sn)
	if err != nil {
		return err
	}

	return m.IOWrite(index, status)
}

func (s *Service) getIO(sn string) (IO, error) {
	m := s.ios[sn]
	if m == nil {
		return nil, errors.New("not found")
	}

	return m, nil
}

// IO模块连接变化
func (s *Service) OnStatus(sn string, status string) {
	s.diag.Debug(fmt.Sprintf("sn:%s status:%s", sn, status))

	ioStatus := []device.Status{
		{
			SN:     sn,
			Type:   device.BaseDeviceTypeIO,
			Status: status,
		},
	}

	s.doDispatch(dispatcherbus.DispatcherDeviceStatus, ioStatus)
}

func (s *Service) OnRecv(string, string) {
	s.diag.Error("OnRecv", errors.New("io Service Not Support OnRecv"))
}

// IO输入输出状态发生变化
func (s *Service) OnChangeIOStatus(sn string, t string, status string) {
	s.diag.Debug(fmt.Sprintf("sn:%s type:%s status:%s", sn, t, status))

	ioContact := IoContact{
		Src: device.BaseDeviceTypeIO,
		SN:  sn,
	}

	if t == IoTypeInput {
		ioContact.Inputs = status
	} else {
		ioContact.Outputs = status
	}

	// IO数据输出状态分发
	s.doDispatch(dispatcherbus.DispatcherIO, ioContact)
}

func (s *Service) doDispatch(name string, data interface{}) {
	if err := s.dispatcherBus.Dispatch(name, data); err != nil {
		s.diag.Error("Dispatch Failed", err)
	}
}
