package storage

import (
	"fmt"
	"github.com/masami10/rush/toml"
	"time"
)

type Config struct {
	Enable       bool          `yaml:"enable"`
	Url          string        `yaml:"db_url"`
	DBName       string        `yaml:"db_name"`
	User         string        `yaml:"db_user"`
	Password     string        `yaml:"db_pwd"`
	MaxConnects  int           `yaml:"max_connection"`
	VacuumPeriod toml.Duration `yaml:"vacuum_period"`
	DataKeep     toml.Duration `yaml:"data_keep"`
}

func NewConfig() Config {
	return Config{
		Enable:       true,
		Url:          "127.0.0.1:5432",
		DBName:       "dbname",
		User:         "user",
		Password:     "pwd",
		MaxConnects:  60,
		VacuumPeriod: toml.Duration(time.Duration(7 * 24 * time.Hour)),
		DataKeep:     toml.Duration(time.Duration(3 * 30 * 24 * time.Hour)),
	}
}

func (c Config) Validate() error {

	if c.VacuumPeriod < toml.Duration(time.Duration(5*24*time.Hour)) {
		return fmt.Errorf("vacuum period %s is less than 5 days", c.VacuumPeriod.String())
	}

	return nil
}
