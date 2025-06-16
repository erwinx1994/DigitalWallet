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
		Server: Server{
			Protocol: "http",
			URL:      "localhost",
			Port:     "1120",
		},
		RequestTimeout: 10, // s
	}
	if !reflect.DeepEqual(*config, expected_config) {
		t.Fatal("Expected: ", expected_config, " Got: ", *config)
	}

}
