package implementation

import (
	"context"
	"database/sql"
	"deposit_service/config"
	"encoding/json"
	"log"
	"reflect"
	"shared/messages"
	"shared/responses"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func Test_DepositService(t *testing.T) {

	message_id := int64(283729)
	currency := "SGD"
	wallet_id := "deposit_service_unit_test"
	deposit_1_amount := "101.11"
	expected_balance_1 := deposit_1_amount
	deposit_2_amount := "10"
	expected_balance_2 := "111.11"
	deposit_3_amount := "100"
	other_currency := "HKD"

	// Load default configuration file
	config_file_path := "../config.yml"
	config, err := config.Load(config_file_path)
	if err != nil {
		t.Fatal("Unable to load configuration file at ", config_file_path)
	}

	// Modify the table and message queue names
	config.RequestsQueue.QueueName = "deposit_requests_queue_test"
	config.ResponsesQueue.QueueName = "deposit_responses_queue_test"
	config.WalletDatabase.BalanceTable = "postgres.test_deposit_service.balances"
	config.WalletDatabase.TransactionsTable = "postgres.test_deposit_service.transactions"

	// Start running balance service
	service := CreateDepositService(config)
	service.Run()
	defer service.Shutdown()

	// Prepare requests queue
	background_context := context.Background()
	requests_queue := redis.NewClient(config.RequestsQueue.GetRedisOptions())
	timeout_context, cancel := context.WithTimeout(background_context, time.Duration(config.RequestsQueue.Timeout)*time.Second)
	_, err = requests_queue.Ping(timeout_context).Result()
	if err != nil {
		cancel()
		t.Fatal("Could not connect to requests queue.", err)
	}
	cancel()

	// Prepare responses queue
	responses_queue := redis.NewClient(config.ResponsesQueue.GetRedisOptions())
	timeout_context, cancel = context.WithTimeout(background_context, time.Duration(config.ResponsesQueue.Timeout)*time.Second)
	_, err = responses_queue.Ping(timeout_context).Result()
	if err != nil {
		cancel()
		t.Fatal("Could not connect to responses queue.", err)
	}
	cancel()

	// Prepare connection to PostgreSQL database
	db, err := sql.Open("postgres", config.WalletDatabase.GetConnectionString())
	if err != nil {
		t.Fatal("Could not create PostgreSQL database object.", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatal("Could not connect to PostgreSQL database.", err)
	}

	// Ensure balances and transactions table are empty when this unit test is finished
	defer func() {
		_, err = db.Exec("delete from " + config.WalletDatabase.BalanceTable)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("delete from " + config.WalletDatabase.TransactionsTable)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Deposit 1
	{
		// Put message into requests queue
		request_message := messages.POST_Deposit{
			Header: messages.Header{
				MessageID: message_id,
				Action:    messages.Action_deposit,
			},
			WalletID: wallet_id,
			Amount:   deposit_1_amount,
			Currency: currency,
		}
		bytes_to_send, err := json.Marshal(request_message)
		if err != nil {
			t.Fatal("Could not serialise message.", err)
		}
		timeout := time.Duration(config.RequestsQueue.Timeout) * time.Second
		queue_name := config.RequestsQueue.QueueName
		timeout_context, cancel = context.WithTimeout(background_context, timeout)
		_, err = requests_queue.LPush(timeout_context, queue_name, bytes_to_send).Result()
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		cancel()

		// Wait for response from response queue
		timeout = time.Duration(config.ResponsesQueue.Timeout) * time.Second
		queue_name = config.ResponsesQueue.QueueName
		timeout_context, cancel = context.WithTimeout(background_context, timeout)
		string_slice, err := responses_queue.BRPop(timeout_context, timeout, queue_name).Result()
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		cancel()

		// string_slice[0] gives the name of the queue
		// string_slice[1] gives the data retrieved from the queue

		// Deserialise JSON data received
		response_message := responses.Deposit{}
		err = json.Unmarshal([]byte(string_slice[1]), &response_message)
		if err != nil {
			t.Fatal(err)
		}

		// Verify response
		expected_response := responses.Deposit{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
			Status:     responses.Status_successful,
			Currency:   request_message.Currency,
			NewBalance: expected_balance_1,
		}
		if !reflect.DeepEqual(response_message, expected_response) {
			t.Fatal("Expected: ", expected_response, ", Got: ", response_message)
		}
	}

	// Deposit 2
	{
		// Put message into requests queue
		request_message := messages.POST_Deposit{
			Header: messages.Header{
				MessageID: message_id,
				Action:    messages.Action_deposit,
			},
			WalletID: wallet_id,
			Amount:   deposit_2_amount,
			Currency: currency,
		}
		bytes_to_send, err := json.Marshal(request_message)
		if err != nil {
			t.Fatal("Could not serialise message.", err)
		}
		timeout := time.Duration(config.RequestsQueue.Timeout) * time.Second
		queue_name := config.RequestsQueue.QueueName
		timeout_context, cancel = context.WithTimeout(background_context, timeout)
		_, err = requests_queue.LPush(timeout_context, queue_name, bytes_to_send).Result()
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		cancel()

		// Wait for response from response queue
		timeout = time.Duration(config.ResponsesQueue.Timeout) * time.Second
		queue_name = config.ResponsesQueue.QueueName
		timeout_context, cancel = context.WithTimeout(background_context, timeout)
		string_slice, err := responses_queue.BRPop(timeout_context, timeout, queue_name).Result()
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		cancel()

		// string_slice[0] gives the name of the queue
		// string_slice[1] gives the data retrieved from the queue

		// Deserialise JSON data received
		response_message := responses.Deposit{}
		err = json.Unmarshal([]byte(string_slice[1]), &response_message)
		if err != nil {
			t.Fatal(err)
		}

		// Verify response
		expected_response := responses.Deposit{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
			Status:     responses.Status_successful,
			Currency:   request_message.Currency,
			NewBalance: expected_balance_2,
		}
		if !reflect.DeepEqual(response_message, expected_response) {
			t.Fatal("Expected: ", expected_response, ", Got: ", response_message)
		}
	}

	// Deposit 3
	{
		// Put message into requests queue
		request_message := messages.POST_Deposit{
			Header: messages.Header{
				MessageID: message_id,
				Action:    messages.Action_deposit,
			},
			WalletID: wallet_id,
			Amount:   deposit_3_amount,
			Currency: other_currency,
		}
		bytes_to_send, err := json.Marshal(request_message)
		if err != nil {
			t.Fatal("Could not serialise message.", err)
		}
		timeout := time.Duration(config.RequestsQueue.Timeout) * time.Second
		queue_name := config.RequestsQueue.QueueName
		timeout_context, cancel = context.WithTimeout(background_context, timeout)
		_, err = requests_queue.LPush(timeout_context, queue_name, bytes_to_send).Result()
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		cancel()

		// Wait for response from response queue
		timeout = time.Duration(config.ResponsesQueue.Timeout) * time.Second
		queue_name = config.ResponsesQueue.QueueName
		timeout_context, cancel = context.WithTimeout(background_context, timeout)
		string_slice, err := responses_queue.BRPop(timeout_context, timeout, queue_name).Result()
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		cancel()

		// string_slice[0] gives the name of the queue
		// string_slice[1] gives the data retrieved from the queue

		// Deserialise JSON data received
		response_message := responses.Deposit{}
		err = json.Unmarshal([]byte(string_slice[1]), &response_message)
		if err != nil {
			t.Fatal(err)
		}

		// Verify response
		expected_response := responses.Deposit{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
			Status:       responses.Status_failed,
			ErrorMessage: "Currency of deposit does not match currency of wallet",
			Currency:     "",
			NewBalance:   "",
		}
		if !reflect.DeepEqual(response_message, expected_response) {
			t.Fatal("Expected: ", expected_response, ", Got: ", response_message)
		}
	}

}
