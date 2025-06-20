package config

import (
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type HTTPServer struct {
	ListenPort      string `yaml:"listen_port"`
	ReadTimeout     int    `yaml:"read_timeout"`     // s
	WriteTimeout    int    `yaml:"write_timeout"`    // s
	IdleTimeout     int    `yaml:"idle_timeout"`     // s
	RetryInterval   int    `yaml:"retry_interval"`   // s
	ShutdownTimeout int    `yaml:"shutdown_timeout"` // s
}

type RedisMessageQueue struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	QueueName string `yaml:"queue_name"`
	Timeout   int    `yaml:"timeout"` // s
}

type Service struct {
	RequestsQueue    RedisMessageQueue `yaml:"redis_requests_queue"`
	ResponsesQueue   RedisMessageQueue `yaml:"redis_responses_queue"`
	CacheWaitTimeout int               `yaml:"cache_wait_timeout"` // s
}

type Config struct {
	HTTPServer                HTTPServer `yaml:"http_server"`
	DepositsService           Service    `yaml:"deposits_service"`
	WithdrawalService         Service    `yaml:"withdrawal_service"`
	TransferService           Service    `yaml:"transfer_service"`
	BalanceService            Service    `yaml:"balance_service"`
	TransactionHistoryService Service    `yaml:"transaction_history_service"`
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
