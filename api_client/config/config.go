package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Protocol string `yaml:"protocol"`
	URL      string `yaml:"url"`
	Port     string `yaml:"port"`
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

func (server *Server) GetURL() string {
	var url string = ""
	if len(server.Port) > 0 {
		url = server.Protocol + "://" + server.URL + ":" + server.Port
	} else {
		url = server.Protocol + "://" + server.URL
	}
	return url
}
