package diagnostic

import (
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	File   string        `yaml:"file"`
	Level  string        `yaml:"level"`
	MaxAge toml.Duration `yaml:"max_age"`
	Rotate toml.Duration `yaml:"rotate"`
}

func NewConfig() Config {
	return Config{
		File:   "STDERR",
		Level:  "DEBUG",
		MaxAge: toml.Duration(time.Duration(7 * 24 * time.Hour)),
		Rotate: toml.Duration(time.Duration(24 * time.Hour)),
	}
}
