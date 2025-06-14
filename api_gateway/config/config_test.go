package config

import (
	"reflect"
	"testing"
)

func Test_LoadConfig(t *testing.T) {

	test_file_path := "../config.yml"
	config, err := Load(test_file_path)
	if err != nil {
		t.Fatal(err)
	}

	expected_config := Config{
		HTTPServer: HTTPServer{
			ListenPort:      "1120",
			ReadTimeout:     60,
			WriteTimeout:    60,
			IdleTimeout:     60,
			RetryInterval:   60,
			ShutdownTimeout: 60,
		},
		DepositsService: Service{
			RequestsQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "deposit_requests_queue",
				Timeout:   5,
			},
			ResponsesQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "deposit_responses_queue",
				Timeout:   5,
			},
		},
		WithdrawalService: Service{
			RequestsQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "withdrawal_requests_queue",
				Timeout:   5,
			},
			ResponsesQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "withdrawal_responses_queue",
				Timeout:   5,
			},
		},
		TransferService: Service{
			RequestsQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "transfer_requests_queue",
				Timeout:   5,
			},
			ResponsesQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "transfer_responses_queue",
				Timeout:   5,
			},
		},
		BalanceService: Service{
			RequestsQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "balance_requests_queue",
				Timeout:   5,
			},
			ResponsesQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "balance_responses_queue",
				Timeout:   5,
			},
		},
		TransactionHistoryService: Service{
			RequestsQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "transaction_history_requests_queue",
				Timeout:   5,
			},
			ResponsesQueue: RedisMessageQueue{
				Host:      "localhost",
				Port:      "1640",
				Username:  "default",
				Password:  "",
				QueueName: "transaction_history_responses_queue",
				Timeout:   5,
			},
		},
	}
	if !reflect.DeepEqual(*config, expected_config) {
		t.Fatal("Expected: ", expected_config, " Got: ", *config)
	}

}
