package config

import (
	"os"

	shared_config "shared/config"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RequestsQueue  shared_config.RedisMessageQueue  `yaml:"redis_requests_queue"`
	ResponsesQueue shared_config.RedisMessageQueue  `yaml:"redis_responses_queue"`
	WalletDatabase shared_config.PostgreSQLDatabase `yaml:"postgresql_wallet_database"`
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
