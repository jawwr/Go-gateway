package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
	Prefix string `yaml:"prefix"`
	Url    string `yaml:"url"`
}

func ReadConfigYAML(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
