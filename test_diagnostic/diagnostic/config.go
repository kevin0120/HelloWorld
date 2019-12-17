package diagnostic

import (
	"time"
)

type Config struct {
	File   string        `yaml:"file"`
	Level  string        `yaml:"level"`
	MaxAge time.Duration `yaml:"max_age"`
	Rotate time.Duration `yaml:"rotate"`
}

func NewConfig() Config {
	return Config{
		File:   "STDERR",
		Level:  "DEBUG",
		MaxAge: time.Duration(31 * 24 * time.Hour),
		Rotate: time.Duration(24 * time.Hour),
	}
}
