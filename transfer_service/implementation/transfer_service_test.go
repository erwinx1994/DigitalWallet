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
	"transfer_service/config"

	"github.com/redis/go-redis/v9"
)

func Test_TransferService(t *testing.T) {

	message_id := int64(283729)
	currency := "SGD"

	source_wallet_id := "source_wallet_id"
	source_initial_balance := 10000

	destination_wallet_id := "destination_wallet_id"
	destination_initial_balance := 10000

	transfer_amount := "10.00"
	source_final_balance_str := "90.00"
	source_final_balance := 9000
	destination_final_balance := 11000

	// Load default configuration file
	config_file_path := "../config.yml"
	config, err := config.Load(config_file_path)
	if err != nil {
		t.Fatal("Unable to load configuration file at ", config_file_path)
	}

	// Modify the table and message queue names
	config.RequestsQueue.QueueName = "transfer_requests_queue_test"
	config.ResponsesQueue.QueueName = "transfer_responses_queue_test"
	config.WalletDatabase.BalanceTable = "postgres.test_transfer_service.balances"
	config.WalletDatabase.TransactionsTable = "postgres.test_transfer_service.transactions"

	// Start running balance service
	service := CreateTransferService(config)
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

	// Create source wallet
	{
		transaction_date_time := time.Now().UTC()
		_, err := db.Exec("insert into "+config.WalletDatabase.TransactionsTable+" (wallet_id, date_and_time, currency, amount) values ($1, $2, $3, $4)", source_wallet_id, transaction_date_time, currency, source_initial_balance)
		if err != nil {
			t.Fatal("Error updating database.", err)
		}

		_, err = db.Exec("insert into "+config.WalletDatabase.BalanceTable+" (wallet_id, currency, balance) values ($1, $2, $3)", source_wallet_id, currency, source_initial_balance)
		if err != nil {
			t.Fatal("Error updating database.", err)
		}
	}

	// Create destination wallet
	{
		transaction_date_time := time.Now().UTC()
		_, err := db.Exec("insert into "+config.WalletDatabase.TransactionsTable+" (wallet_id, date_and_time, currency, amount) values ($1, $2, $3, $4)", destination_wallet_id, transaction_date_time, currency, destination_initial_balance)
		if err != nil {
			t.Fatal("Error updating database.", err)
		}

		_, err = db.Exec("insert into "+config.WalletDatabase.BalanceTable+" (wallet_id, currency, balance) values ($1, $2, $3)", destination_wallet_id, currency, destination_initial_balance)
		if err != nil {
			t.Fatal("Error updating database.", err)
		}
	}

	// Transfer some money from source to destination wallets
	{
		// Put message into requests queue
		request_message := messages.POST_Transfer{
			Header: messages.Header{
				MessageID: message_id,
				Action:    messages.Action_transfer,
			},
			SourceWalletID:      source_wallet_id,
			DestinationWalletID: destination_wallet_id,
			Amount:              transfer_amount,
			Currency:            currency,
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
		response_message := responses.Transfer{}
		err = json.Unmarshal([]byte(string_slice[1]), &response_message)
		if err != nil {
			t.Fatal(err)
		}

		// Verify response
		expected_response := responses.Transfer{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
			Status:     responses.Status_successful,
			Currency:   request_message.Currency,
			NewBalance: source_final_balance_str,
		}
		if !reflect.DeepEqual(response_message, expected_response) {
			t.Fatal("Expected: ", expected_response, ", Got: ", response_message)
		}

		get_balance, err := db.Prepare("select balance from " + service.config.WalletDatabase.BalanceTable + " where wallet_id=$1")
		if err != nil {
			t.Fatal("Unable to prepare SQL statement.")
		}
		defer get_balance.Close()

		// Verify balance of source wallet
		var balance int64 = 0
		err = get_balance.QueryRow(request_message.SourceWalletID).Scan(&balance)
		if err != nil {
			t.Fatal(err)
		}
		if balance != int64(source_final_balance) {
			t.Error("Expected: ", source_final_balance, ", Got: ", balance)
		}

		// Verify balance of destination wallet
		err = get_balance.QueryRow(request_message.DestinationWalletID).Scan(&balance)
		if err != nil {
			t.Fatal(err)
		}
		if balance != int64(destination_final_balance) {
			t.Error("Expected: ", destination_final_balance, ", Got: ", balance)
		}
	}

}
