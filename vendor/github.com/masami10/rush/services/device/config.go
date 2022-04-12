package device

type Config struct {
	Enable bool `yaml:"enable"`
}

func NewConfig() Config {

	return Config{
		Enable: true,
	}
}

func (c Config) Validate() error {

	return nil
}
