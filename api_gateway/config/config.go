package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenPort      string `yaml:"listen_port"`
	ReadTimeout     int    `yaml:"read_timeout"`     // s
	WriteTimeout    int    `yaml:"write_timeout"`    // s
	IdleTimeout     int    `yaml:"idle_timeout"`     // s
	RetryInterval   int    `yaml:"retry_interval"`   // s
	ShutdownTimeout int    `yaml:"shutdown_timeout"` // s
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
