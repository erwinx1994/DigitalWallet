package implementation

import (
	"api_gateway/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"shared/messages"
	"testing"
	"time"
)

func Test_APIGateway(t *testing.T) {

	config, err := config.Load("../config.yml")
	if err != nil {
		t.Fatal(err)
	}
	config.BalanceService.RequestsQueue.QueueName = "api_gateway_requests_queue_test"
	config.BalanceService.ResponsesQueue.QueueName = "api_gateway_responses_queue_test"

	// Start running test service
	test_service_ := create_test_service(&config.BalanceService)
	test_service_.run()
	defer test_service_.shutdown()

	// Start running API gateway
	api_gateway, err := CreateAPIGateway(config)
	if err != nil {
		t.Fatal(err)
	}
	api_gateway.Run()
	defer api_gateway.Shutdown()

	time.Sleep(2 * time.Second)

	// Create HTTP client
	http_timeout := 5 // s
	http_client := http.Client{
		Timeout: time.Duration(http_timeout) * time.Second,
	}
	url := "http://localhost:1120/test"

	// Test HTTP GET request
	{
		response, err := http_client.Get(url)
		if err != nil {
			t.Fatal(err.Error())
		}

		// Check that response is as expected
		length_of_body := response.ContentLength
		bytes := make([]byte, length_of_body)
		_, err = response.Body.Read(bytes)
		if err != nil && !errors.Is(err, io.EOF) {
			response.Body.Close()
			t.Fatal(err)
		}
		response.Body.Close()

		expected_data := []byte("{\"header\":{\"id\":1,\"action\":4},\"status\":1,\"error_message\":\"Test service: I received your message: \"}")
		if !reflect.DeepEqual(bytes, expected_data) {
			t.Error("Expected: ", string(expected_data), ", Got: ", string(bytes))
		}
	}

	// Test HTTP POST request
	{
		// Prepare body of HTTP POST message
		body := messages.GET_Balance{
			WalletID: "This is test message 2.",
		}
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		// Send HTTP POST request and get response
		response, err := http_client.Post(url, "application/json", bytes.NewReader(data))
		if err != nil {
			t.Fatal(err.Error())
		}

		// Check that response is as expected
		length_of_body := response.ContentLength
		bytes := make([]byte, length_of_body)
		_, err = response.Body.Read(bytes)
		if err != nil && !errors.Is(err, io.EOF) {
			response.Body.Close()
			t.Fatal(err)
		}
		response.Body.Close()

		expected_data := []byte("{\"header\":{\"id\":2,\"action\":4},\"status\":1,\"error_message\":\"Test service: I received your message: This is test message 2.\"}")
		if !reflect.DeepEqual(bytes, expected_data) {
			t.Fatal("Expected: ", string(expected_data), ", Got: ", string(bytes))
		}
	}
}
