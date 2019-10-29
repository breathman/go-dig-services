package config

import (
	"io/ioutil"
	"log" // nolint
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config interface {
	GetConfigValue(value string) (interface{}, error)
}

type Main struct {
	Logger Logger `yaml:"logger"`
}

type Logger struct {
	LoggerLevel string `yaml:"log_level"`
	IsDevMode   bool   `yaml:"is_dev_mode"`
}

type Influx struct {
	Host   string `yaml:"host" validate:"required"`
	Port   string `yaml:"port" validate:"required"`
	DBName string `yaml:"db_name" validate:"required"`
}

func ReadConfig(configEnvName string, config interface{}) error {
	configPath := os.Getenv(configEnvName)
	if configPath == `` {
		log.Fatalf(`%q env must be defined`, configEnvName)
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.Wrap(err, `failed to read config file`)
	}

	if err = yaml.Unmarshal(configBytes, config); err != nil {
		return errors.Wrap(err, `failed to unmarshal yaml config`)
	}
	return nil
}
