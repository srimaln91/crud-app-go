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
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Name               string `yaml:"name"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	PoolSize           int    `yaml:"pool_size"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
	ConnMaxLifeTime    int    `yaml:"conn_max_lifetime"`
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
