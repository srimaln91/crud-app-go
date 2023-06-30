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
	Name string `yaml:"name"`
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
