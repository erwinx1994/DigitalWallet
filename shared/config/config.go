package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisMessageQueue struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	QueueName string `yaml:"queue_name"`
	Timeout   int    `yaml:"timeout"`
}

type PostgreSQLDatabase struct {
	Host              string `yaml:"host"`
	Port              string `yaml:"port"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	Database          string `yaml:"database"`
	BalanceTable      string `yaml:"balance_table"`
	TransactionsTable string `yaml:"transactions_table"`
}

func (message_queue *RedisMessageQueue) GetRedisOptions() *redis.Options {
	options := &redis.Options{
		Addr:     message_queue.Host + ":" + message_queue.Port,
		Username: message_queue.Username,
		Password: message_queue.Password,
	}
	return options
}

func (postgres *PostgreSQLDatabase) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgres.Host, postgres.Port, postgres.Username, postgres.Password, postgres.Database)
}
