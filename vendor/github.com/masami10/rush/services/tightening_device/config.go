package tightening_device

const (
	IOSelectorType = "io"
)

type TighteningDeviceConfig struct {
	// 控制器型号
	Model string `yaml:"model" json:"model"`

	// 控制器协议类型
	Protocol string `yaml:"protocol" json:"protocol"`

	// 连接地址(如果在控制器上配了连接地址，则下属所有工具共用此地址进行通信)
	Endpoint string `yaml:"endpoint" json:"endpoint"`

	// 控制器序列号
	SN string `yaml:"sn" json:"sn"`

	// 控制器名字
	ControllerName string `yaml:"name" json:"name"`

	// 工具列表
	Tools []ToolConfig `yaml:"tools" json:"children"`
}

type ToolConfig struct {
	// 工具序列号
	SN string `yaml:"sn" json:"sn"`

	// 工具通道号
	Channel int `yaml:"channel" json:"channel"`

	// 连接地址
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

type SocketSelectorConfig struct {
	Enable   bool   `yaml:"enable"`
	Endpoint string `yaml:"endpoint"`
	Type     string `yaml:"type"`
}

type Config struct {
	Enable         bool                     `yaml:"enable"`
	SocketSelector SocketSelectorConfig     `yaml:"socket_selector"`
	Devices        []TighteningDeviceConfig `yaml:"devices"`
}

func NewConfig() Config {

	return Config{
		Enable:  true,
		Devices: []TighteningDeviceConfig{},
	}
}

func (c Config) Validate() error {

	return nil
}
