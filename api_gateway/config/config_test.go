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
		ListenPort:      "1120",
		ReadTimeout:     60,
		WriteTimeout:    60,
		IdleTimeout:     60,
		RetryInterval:   60,
		ShutdownTimeout: 60,
	}
	if !reflect.DeepEqual(*config, expected_config) {
		t.Fatal("Expected: ", expected_config, " Got: ", *config)
	}

}
