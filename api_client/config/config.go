package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	URL  string `yaml:"url"`
	Port string `yaml:"port"`
}

type Config struct {
	Server         Server `yaml:"server"`
	RequestTimeout int    `yaml:"request_timeout"`
}

func Load(filepath string) (*Config, error) {

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config := Config{}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
