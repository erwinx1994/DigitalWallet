package config

import (
	"reflect"
	"testing"
)

func Test_LoadConfig(t *testing.T) {

	test_file_path := "config_test.yml"
	config, err := Load(test_file_path)
	if err != nil {
		t.Fatal(err)
	}

	expected_config := Config{
		ListenPort:      "1120",
		ReadTimeout:     5,
		WriteTimeout:    5,
		IdleTimeout:     5,
		RetryInterval:   5,
		ShutdownTimeout: 5,
	}
	if !reflect.DeepEqual(*config, expected_config) {
		t.Fatal("Expected: ", expected_config, " Got: ", *config)
	}

}
