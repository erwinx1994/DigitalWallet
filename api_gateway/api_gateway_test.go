package main

import (
	"api_gateway/config"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_APIGateway(t *testing.T) {

	config, err := config.Load("./config/config_test.yml")
	if err != nil {
		t.Fatal(err)
	}

	// Start running API gateway
	api_gateway := &APIGateway{
		config:      config,
		http_server: nil,
	}
	api_gateway.run()
	defer api_gateway.shutdown()

	time.Sleep(2 * time.Second)

	// Make test HTTP GET request
	http_timeout := 5 // s
	http_client := http.Client{
		Timeout: time.Duration(http_timeout) * time.Second,
	}
	get_balance := http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:1120",
			Path:   "/test",
		},
	}

	// Get response from server
	response, err := http_client.Do(&get_balance)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer response.Body.Close()

	// Check that response is as expected
	length_of_body := response.ContentLength
	bytes := make([]byte, length_of_body)
	_, err = response.Body.Read(bytes)
	if err != nil && !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}
	expected_data := []byte("This is a test response")
	if !reflect.DeepEqual(bytes, expected_data) {
		t.Fatal("Expected: ", string(expected_data), ", Got: ", string(bytes))
	}

}
