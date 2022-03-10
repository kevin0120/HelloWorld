package device

//status
const (
	BaseDeviceStatusOnline  = "online"
	BaseDeviceStatusOffline = "offline"
	//拧紧工具特有的
	BaseDeviceStatusRunning   = "running"
	BaseDeviceStatusUnRunning = "stopping"
	BaseDeviceStatusEnabled   = "enabled"
	BaseDeviceStatusDisabled  = "disabled"
)

// common device type Define
const (
	BaseDeviceTypeIO = "io"
	//BaseDeviceTypeReader     = "reader"
	BaseDeviceTighteningTool   = "tightening_tool"
	BaseDeviceTypeScanner      = "scanner"
	BaseDeviceTypeMeasureRuler = "measure_ruler"
)
