package device

import (
	"github.com/masami10/rush/services/dispatcherbus"
	"github.com/masami10/rush/utils"
)

type IParentService interface {
	OnStatus(string, string)
	OnRecv(string, string)
}

type IBaseDevice interface {
	GenerateDispatcherNameBySerialNumber(base string) string

	GetParentService() IParentService

	WillStart() error //即将启动

	Start() error

	Stop() error
	// 设备状态
	Status() string

	OnDeviceStatus(string) //状态发生变化

	DoOnDeviceStatus(symbol string, status string) error //执行相应业务

	OnDeviceRecv(string) error //接收到数据

	DoOnDeviceRecv(symbol string, status string) error //执行相应业务

	// 设备类型
	DeviceType() string

	// 子设备
	Children() map[string]IBaseDevice

	// 设备配置
	Config() interface{}

	// 设备运行数据
	Data() interface{}

	//设备序列号唯一追踪号
	SerialNumber() string

	SetSerialNumber(string)
}

type Dispatcher interface {
	Create(name string, len int) error
	Start(name string) error
	Dispatch(name string, data interface{}) error
	LaunchDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
	Release(name string, handler string) error
	Register(name string, handler *utils.DispatchHandlerStruct)
	ReleaseDispatchersByHandlerMap(dispatcherMap dispatcherbus.DispatcherMap)
}

type IDeviceService interface {
	AddDevice(sn string, d IBaseDevice)
}
