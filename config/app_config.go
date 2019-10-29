package config

import (
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

type APPConfig struct {
	Main   `yaml:"main"`
	Influx `yaml:"influx"`
}

func NewAPPConfig() (*APPConfig, error) {
	var c APPConfig
	err := ReadConfig(`CONFIG_PATH`, &c)
	if err != nil {
		return nil, errors.Wrap(err, `failed to read config file`)
	}

	return &c, validator.New().Struct(c)
}
