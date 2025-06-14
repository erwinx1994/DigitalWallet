package config

import (
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type RedisMessageQueue struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	QueueName string `yaml:"queue_name"`
	Timeout   int    `yaml:"timeout"`
}

type Config struct {
	RequestsQueue  RedisMessageQueue `yaml:"redis_requests_queue"`
	ResponsesQueue RedisMessageQueue `yaml:"redis_responses_queue"`
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

func (message_queue *RedisMessageQueue) GetRedisOptions() *redis.Options {
	options := &redis.Options{
		Addr:     message_queue.Host + ":" + message_queue.Port,
		Username: message_queue.Username,
		Password: message_queue.Password,
	}
	return options
}
