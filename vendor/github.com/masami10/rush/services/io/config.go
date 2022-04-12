package io

import (
	"github.com/masami10/rush/toml"
	"github.com/pkg/errors"
	"time"
)

type ConfigIO struct {
	SN      string        `yaml:"sn"`
	Model   string        `yaml:"model"`
	Speed   uint          `yaml:"speed"` // 串口下为波特率
	Address string        `yaml:"address"`
	ReadItv toml.Duration `yaml:"read_itv"`
}

type Config struct {
	Enable       bool          `yaml:"enable"`
	FlashInteval toml.Duration `yaml:"flash_inteval"`
	IOS          []ConfigIO    `yaml:"ios"`
}

func NewConfig() Config {
	return Config{
		Enable:       true,
		FlashInteval: toml.Duration(time.Second * 1),
		IOS: []ConfigIO{
			{
				SN:      "1",
				Model:   ModelMoxaE1212,
				Address: "modbustcp://127.0.0.1:502",
				ReadItv: toml.Duration(300 * time.Millisecond),
			},
			{
				SN:      "2",
				Model:   ModelMoxaE1212,
				Address: "rtu:///dev/ttyUSB0",
				Speed:   19200,
				ReadItv: toml.Duration(300 * time.Millisecond),
			},
		},
	}
}

func (c Config) Validate() error {

	for _, io := range c.IOS {
		_, exist := VendorModels[io.Model]
		if !exist {
			return errors.New("Vendor Not Found")
		}
	}

	return nil
}
