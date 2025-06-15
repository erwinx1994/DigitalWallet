package config

import (
	"reflect"
	shared_config "shared/config"
	"testing"
)

func Test_LoadConfig(t *testing.T) {

	test_file_path := "../config.yml"
	config, err := Load(test_file_path)
	if err != nil {
		t.Fatal(err)
	}

	expected_config := Config{
		RequestsQueue: shared_config.RedisMessageQueue{
			Host:      "localhost",
			Port:      "1640",
			Username:  "default",
			Password:  "",
			QueueName: "balance_requests_queue",
			Timeout:   5,
		},
		ResponsesQueue: shared_config.RedisMessageQueue{
			Host:      "localhost",
			Port:      "1640",
			Username:  "default",
			Password:  "",
			QueueName: "balance_responses_queue",
			Timeout:   5,
		},
		WalletDatabase: shared_config.PostgreSQLDatabase{
			Host:              "localhost",
			Port:              "5432",
			Username:          "postgres",
			Password:          "postgres",
			Database:          "postgres",
			BalanceTable:      "postgres.wallet.balances",
			TransactionsTable: "postgres.wallet.transactions",
		},
	}
	if !reflect.DeepEqual(*config, expected_config) {
		t.Fatal("Expected: ", expected_config, " Got: ", *config)
	}

}
