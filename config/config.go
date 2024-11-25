package config

import (
	"io/ioutil"

	"github.com/srimaln91/crud-app-go/log"
	yaml "gopkg.in/yaml.v3"
)

type AppConfig struct {
	HTTP struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
	Logger struct {
		Level log.Level `yaml:"level"`
	} `yaml:"logger"`
	Database DBConfig `yaml:"database"`
}

type DBConfig struct {
	Database           string `yaml:"database"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	PoolSize           int    `yaml:"pool-size"`
	MaxIdleConnections int    `yaml:"max-idle-conns"`
	ConnMaxLifeTime    int    `yaml:"max-conn-lifetime"`
}

func Parse(path string) (*AppConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(AppConfig)
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
