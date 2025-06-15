package implementation

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"reflect"
	"shared/messages"
	"shared/responses"
	"testing"
	"time"
	"withdraw_service/config"

	"github.com/redis/go-redis/v9"
)

func Test_WithdrawService(t *testing.T) {

	message_id := int64(283729)
	currency := "SGD"
	wallet_id := "withdraw_service_unit_test"
	initial_amount := int64(10000)
	withdraw_1_amount := "10.00"
	expected_balance_1 := "90.00"

	// Load default configuration file
	config_file_path := "../config.yml"
	config, err := config.Load(config_file_path)
	if err != nil {
		t.Fatal("Unable to load configuration file at ", config_file_path)
	}

	// Modify the table and message queue names
	config.RequestsQueue.QueueName = "withdrawal_requests_queue_test"
	config.ResponsesQueue.QueueName = "withdrawal_responses_queue_test"
	config.WalletDatabase.BalanceTable = "postgres.test_withdraw_service.balances"
	config.WalletDatabase.TransactionsTable = "postgres.test_withdraw_service.transactions"

	// Start running balance service
	service := CreateWithdrawService(config)
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

	// Create a new wallet with an initial balance
	{
		transaction_date_time := time.Now().UTC()
		_, err := db.Exec("insert into "+config.WalletDatabase.TransactionsTable+" (wallet_id, date_and_time, currency, amount) values ($1, $2, $3, $4)", wallet_id, transaction_date_time, currency, initial_amount)
		if err != nil {
			t.Fatal("Error updating database.", err)
		}

		_, err = db.Exec("insert into "+config.WalletDatabase.BalanceTable+" (wallet_id, currency, balance) values ($1, $2, $3)", wallet_id, currency, initial_amount)
		if err != nil {
			t.Fatal("Error updating database.", err)
		}
	}

	// Withdraw some money from wallet
	{
		// Put message into requests queue
		request_message := messages.POST_Withdraw{
			Header: messages.Header{
				MessageID: message_id,
				Action:    messages.Action_withdraw,
			},
			WalletID: wallet_id,
			Amount:   withdraw_1_amount,
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
		response_message := responses.Withdraw{}
		err = json.Unmarshal([]byte(string_slice[1]), &response_message)
		if err != nil {
			t.Fatal(err)
		}

		// Verify response
		expected_response := responses.Withdraw{
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

}
