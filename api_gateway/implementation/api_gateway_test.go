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

	// Start running API gateway
	api_gateway := &APIGateway{
		config:      config,
		http_server: nil,
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

		expected_data := []byte("{\"Message\":\"Hi. This is a test response.\"}")
		if !reflect.DeepEqual(bytes, expected_data) {
			t.Fatal("Expected: ", string(expected_data), ", Got: ", string(bytes))
		}
	}

	// Test HTTP POST request
	{
		// Prepare body of HTTP POST message
		body := messages.POST_Deposit{
			Amount:   "100.20",
			Currency: "SGD",
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

		expected_data := []byte("{\"amount\":\"100.20\",\"currency\":\"SGD\"}")
		if !reflect.DeepEqual(bytes, expected_data) {
			t.Fatal("Expected: ", string(expected_data), ", Got: ", string(bytes))
		}
	}

}
