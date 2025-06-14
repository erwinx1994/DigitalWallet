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
	}
	if !reflect.DeepEqual(*config, expected_config) {
		t.Fatal("Expected: ", expected_config, " Got: ", *config)
	}

}
