package implementation

import (
	"balance_service/config"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"shared/messages"
	"shared/responses"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
Ensure that an instance of Redis is running for this test.
*/

func Test_BalanceService(t *testing.T) {

	message_id := int64(283729)
	currency := "SGD"
	wallet_id := "balance_service_unit_test"
	balance := int64(10000)

	// Load default configuration file
	config_file_path := "../config.yml"
	config, err := config.Load(config_file_path)
	if err != nil {
		t.Fatal("Unable to load configuration file at ", config_file_path)
	}

	// Modify the table and message queue names
	config.RequestsQueue.QueueName = "balance_requests_queue_test"
	config.ResponsesQueue.QueueName = "balance_responses_queue_test"
	config.WalletDatabase.BalanceTable = "postgres.test_wallet.balances"
	config.WalletDatabase.TransactionsTable = "postgres.test_wallet.transactions"

	// Start running balance service
	service := CreateBalanceService(config)
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

	// Add test row to balances table
	_, err = db.Exec("insert into "+config.WalletDatabase.BalanceTable+" (wallet_id, currency, balance) values ($1,$2,$3)", wallet_id, currency, balance)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure test row gets deleted when unit test is finished
	defer func() {
		_, err = db.Exec("delete from "+config.WalletDatabase.BalanceTable+" where wallet_id=$1", wallet_id)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Put message into requests queue
	request_message := messages.GET_Balance{
		Header: messages.Header{
			MessageID: message_id,
			Action:    messages.Action_get_balance,
		},
		WalletID: wallet_id,
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
	response_message := responses.Balance{}
	err = json.Unmarshal([]byte(string_slice[1]), &response_message)
	if err != nil {
		t.Fatal(err)
	}

	// Verify response
	if response_message.Header.MessageID != request_message.Header.MessageID {
		t.Error("Expected: ", request_message.Header.MessageID, ", Got: ", response_message.Header.MessageID)
	}
	if response_message.Header.Action != request_message.Header.Action {
		t.Error("Expected: ", request_message.Header.Action, ", Got: ", response_message.Header.Action)
	}

	// Log responses
	message := ""
	switch response_message.Status {
	case responses.Status_unknown:
		message += "Unknown response"
	case responses.Status_successful:
		message += "Successful balance request"
	case responses.Status_failed:
		message += "Failed balance request"
	}
	message += ":" + response_message.Currency + " " + response_message.Balance
	t.Log(message)

}
