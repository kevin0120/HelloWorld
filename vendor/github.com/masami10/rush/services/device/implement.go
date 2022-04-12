package device

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"sync"
)

type BaseDeviceDiagnostic interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

type Notify interface {
	NotifyAll() error
	NotifySomeOne(symbol string) error
}

func CreateBaseDevice(deviceType string, d BaseDeviceDiagnostic, service IParentService, serialNumber string) BaseDevice {
	c := BaseDevice{
		children:     &sync.Map{},
		diag:         d,
		service:      service,
		status:       atomic.NewString(BaseDeviceStatusOffline),
		deviceType:   deviceType,
		serialNumber: serialNumber,
	}

	return c
}

type BaseDevice struct {
	status     *atomic.String
	deviceType string
	diag       Diagnostic
	//mtxChildren  sync.Mutex
	service      IParentService
	children     *sync.Map
	Cfg          interface{}
	manufacture  string // 设备厂商名
	serialNumber string
}

func (s *BaseDevice) GenerateDispatcherNameBySerialNumber(base string) string {
	return fmt.Sprintf("%s@%s@%s@%s", base, s.serialNumber, s.deviceType, s.manufacture)
}

func (s *BaseDevice) Manufacture() string {
	return s.manufacture
}

func (s *BaseDevice) SetManufacture(m string) {
	s.manufacture = m
}

func (s *BaseDevice) DeviceType() string {
	return s.deviceType
}

func (s *BaseDevice) GetParentService() IParentService {
	return s.service
}

func (s *BaseDevice) WillStart() error {
	s.UpdateStatus(BaseDeviceStatusOffline)
	return nil
}

func (s *BaseDevice) Start() error {
	return s.WillStart()
}

func (s *BaseDevice) Stop() error {
	s.UpdateStatus(BaseDeviceStatusOffline)
	return nil
}

func (s *BaseDevice) DoOnDeviceStatus(symbol string, status string) error {
	if s.service == nil {
		return errors.New("Please Inject Device Parent Service First")
	}
	s.service.OnStatus(symbol, status)
	return nil
}

func (s *BaseDevice) OnDeviceStatus(status string) {
	if status != BaseDeviceStatusOnline && status != BaseDeviceStatusOffline {
		s.diag.Error("onStatus", errors.Errorf("Status: %s Is Not Support", status))
	}
	ss := s.Status()
	if ss == status {
		return
	}
	s.UpdateStatus(status)
	if err := s.DoOnDeviceStatus(s.serialNumber, status); err != nil {
		s.diag.Error("onStatus DoOnDeviceStatus", err)
	}
}

func (s *BaseDevice) DoOnDeviceRecv(symbol, msg string) error {
	if s.service == nil {
		return errors.New("Please Inject Device Parent Service First")
	}
	s.service.OnRecv(symbol, msg)
	return nil
}

func (s *BaseDevice) OnDeviceRecv(msg string) error {
	return s.DoOnDeviceRecv(s.serialNumber, msg)
}

func (s *BaseDevice) SerialNumber() string {
	return s.serialNumber
}

func (s *BaseDevice) SetSerialNumber(serialNumber string) {
	s.serialNumber = serialNumber
}

func (s *BaseDevice) UpdateStatus(status string) {
	s.status.Store(status)
}

func (s *BaseDevice) Status() string {
	return s.status.Load()
}

func (s *BaseDevice) AddChildren(sn string, device IBaseDevice) {
	s.children.Store(sn, device)
}

func (s *BaseDevice) ChildrenRange(f func(key string, value IBaseDevice) bool) {
	fn := func(key, value interface{}) bool {
		sn := key.(string)
		v := value.(IBaseDevice)
		return f(sn, v)
	}
	s.children.Range(fn)
}

func (s *BaseDevice) FindChild(key string) IBaseDevice {
	if v, ok := s.children.Load(key); ok {
		return v.(IBaseDevice)
	}
	return nil
}

func (s *BaseDevice) Children() map[string]IBaseDevice {
	var ret = make(map[string]IBaseDevice)
	fn := func(key, value interface{}) bool {
		sn := key.(string)
		v := value.(IBaseDevice)
		ret[sn] = v
		return true
	}
	s.children.Range(fn)
	return ret
}

// 设备配置
func (s *BaseDevice) Config() interface{} {
	return s.Cfg
}

// 设备运行数据
func (s *BaseDevice) Data() interface{} {
	return nil
}
