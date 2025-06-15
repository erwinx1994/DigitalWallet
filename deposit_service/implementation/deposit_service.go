package implementation

import (
	"context"
	"database/sql"
	"deposit_service/config"
	"encoding/json"
	"errors"
	"log"
	"shared/messages"
	"shared/responses"
	"shared/utilities"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type DepositService struct {
	config             *config.Config
	is_alive           atomic.Bool
	waitgroup          sync.WaitGroup
	background_context context.Context
	requests_queue     *redis.Client
	responses_queue    *redis.Client
}

func CreateDepositService(config *config.Config) *DepositService {
	service := &DepositService{
		config:             config,
		background_context: context.Background(),
	}
	return service
}

func (service *DepositService) prepare_redis_clients() error {

	{
		// Prepare requests queue
		requests_queue := redis.NewClient(service.config.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(service.background_context, time.Duration(service.config.RequestsQueue.Timeout)*time.Second)
		_, err := requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()

		service.requests_queue = requests_queue
	}

	{
		// Prepare responses queue
		responses_queue := redis.NewClient(service.config.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(service.background_context, time.Duration(service.config.ResponsesQueue.Timeout)*time.Second)
		_, err := responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()

		service.responses_queue = responses_queue
	}

	return nil
}

func (service *DepositService) send_response(response_message *responses.Deposit) {
	bytes_to_send, err := json.Marshal(response_message)
	if err != nil {
		log.Println("Failed to serialise response message. Should not happen in production.")
		// In practice, we will need an error notification system. I have skipped
		// building an error notification system due to time constraints.
		return
	}

	// Put response into responses queue
	timeout := time.Duration(service.config.ResponsesQueue.Timeout) * time.Second
	queue_name := service.config.ResponsesQueue.QueueName
	timeout_context, cancel := context.WithTimeout(service.background_context, timeout)
	_, err = service.responses_queue.LPush(timeout_context, queue_name, bytes_to_send).Result()
	if err != nil {
		cancel()
		log.Println("Failed to put response into responses queue.")
		return
	}
	cancel()
}

func (service *DepositService) send_failed_response(request_message *messages.POST_Deposit) {
	response_message := responses.Deposit{
		Header: responses.Header{
			MessageID: request_message.Header.MessageID,
			Action:    request_message.Header.Action,
		},
		Status: responses.Status_failed,
	}
	service.send_response(&response_message)
}

func (service *DepositService) async_run() {
	defer func() {
		service.waitgroup.Done()
		log.Println("Shutdown deposit service.")
	}()

	log.Println("Started up deposit service.")

	err := service.prepare_redis_clients()
	if err != nil {
		log.Fatal("Could not connect to Redis server: ", err)
	}
	log.Println("Created clients for Redis message queues.")

	// Open connection to PostgreSQL database
	db, err := sql.Open("postgres", service.config.WalletDatabase.GetConnectionString())
	if err != nil {
		log.Fatal("Could not create PostgreSQL database object.", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to PostgreSQL database.", err)
	}
	log.Println("Connected to PostgreSQL database.")

	// Prepare commonly used SQL statements
	get_currency_balance, err := db.Prepare("select currency, balance from " + service.config.WalletDatabase.BalanceTable + " where wallet_id=$1")
	if err != nil {
		log.Fatal("Unable to prepare SQL statement.")
	}
	defer get_currency_balance.Close()

	insert_transaction, err := db.Prepare("insert into " + service.config.WalletDatabase.TransactionsTable + " (wallet_id, date_and_time, currency, amount) values ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal("Unable to prepare SQL statement.")
	}
	defer insert_transaction.Close()

	update_balance, err := db.Prepare("update " + service.config.WalletDatabase.BalanceTable + " set balance=$1 where wallet_id=$2")
	if err != nil {
		log.Fatal("Unable to prepare SQL statement.")
	}
	defer update_balance.Close()

	insert_new_balance, err := db.Prepare("insert into " + service.config.WalletDatabase.BalanceTable + " (wallet_id, currency, balance) values ($1, $2, $3)")
	if err != nil {
		log.Fatal("Unable to prepare SQL statement.")
	}
	defer insert_new_balance.Close()

	// Service continues running until terminated by user
	for service.is_alive.Load() {

		// Get next request from requests queue
		timeout := time.Duration(service.config.RequestsQueue.Timeout) * time.Second
		queue_name := service.config.RequestsQueue.QueueName
		timeout_context, cancel := context.WithTimeout(service.background_context, timeout)
		string_slice, err := service.requests_queue.BRPop(timeout_context, timeout, queue_name).Result()
		if err != nil {
			cancel()
			continue
		}
		cancel()

		// string_slice[0] gives the name of the queue
		// string_slice[1] gives the data retrieved from the queue

		// Deserialise JSON data received
		request_message := messages.POST_Deposit{}
		err = json.Unmarshal([]byte(string_slice[1]), &request_message)
		if err != nil {
			log.Println("Failed to deserialise JSON message. Should not happen in production.")
			// In practice, we will need an error notification system. Building an error
			// notification is skipped due to time constraints.
			continue
		}

		// Verify that the correct message was received
		if request_message.Header.Action != messages.Action_deposit {
			log.Println("Incorrect message received.")
			continue
		}

		// Very that request message is valid
		deposit_amount, err := utilities.Convert_display_to_database_format(request_message.Amount)
		if err != nil {
			log.Println("Invalid amount.", err)
			service.send_failed_response(&request_message)
			continue
		}

		// Query PostgreSQL database
		transaction_date_time := time.Now().UTC()
		db_transaction, err := db.Begin()
		if err != nil {
			log.Println("Error creating transaction.")
			service.send_failed_response(&request_message)
			continue
		}

		// 1. Check if wallet already exists and determine its currency. Return an error
		// if the currency of the wallet does not match the deposit.
		var currency string = ""
		var balance int64 = 0
		tx_get_currency_balance := db_transaction.Stmt(get_currency_balance)
		err = tx_get_currency_balance.QueryRow(request_message.WalletID).Scan(&currency, &balance)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				db_transaction.Rollback()
				log.Println("Error reading from PostgreSQL database: ", err)
				service.send_failed_response(&request_message)
				continue
			}
		}

		if err != nil {
			// Wallet does not exist yet

			// Create a new deposit transaction.
			tx_insert_transaction := db_transaction.Stmt(insert_transaction)
			_, err := tx_insert_transaction.Exec(request_message.WalletID, transaction_date_time, request_message.Currency, deposit_amount)
			if err != nil {
				db_transaction.Rollback()
				log.Println("Database transaction failed.", err)
				service.send_failed_response(&request_message)
				continue
			}

			// Create a new entry in the balance table.
			balance = deposit_amount
			tx_insert_new_balance := db_transaction.Stmt(insert_new_balance)
			_, err = tx_insert_new_balance.Exec(request_message.WalletID, request_message.Currency, balance)
			if err != nil {
				db_transaction.Rollback()
				log.Println("Database transaction failed.", err)
				service.send_failed_response(&request_message)
				continue
			}
		} else {

			if currency == request_message.Currency {
				// Wallet already exists and its currency matches with the currency of the
				// deposit

				// Create a new deposit transaction.
				tx_insert_transaction := db_transaction.Stmt(insert_transaction)
				_, err := tx_insert_transaction.Exec(request_message.WalletID, transaction_date_time, request_message.Currency, deposit_amount)
				if err != nil {
					db_transaction.Rollback()
					log.Println("Database transaction failed.", err)
					service.send_failed_response(&request_message)
					continue
				}

				// Update the balance of the wallet
				balance += deposit_amount
				tx_update_balance := db_transaction.Stmt(update_balance)
				_, err = tx_update_balance.Exec(balance, request_message.WalletID)
				if err != nil {
					db_transaction.Rollback()
					log.Println("Database transaction failed.", err)
					service.send_failed_response(&request_message)
					continue
				}

			} else {
				// Do not proceed with deposit if wallet already exists and its currency does
				// not match with the currency of the deposit
				db_transaction.Rollback()
				log.Println("Mismatching wallet and deposit currencies", err)
				service.send_failed_response(&request_message)
				continue
			}

		}

		// Commit database transaction
		err = db_transaction.Commit()
		if err != nil {
			db_transaction.Rollback()
			log.Println("Database transaction failed.", err)
			service.send_failed_response(&request_message)
			continue
		}

		// Prepare response
		response_message := responses.Deposit{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
			Status:     responses.Status_successful,
			Currency:   request_message.Currency,
			NewBalance: utilities.Convert_database_to_display_format(balance),
		}
		service.send_response(&response_message)

	}
}

func (service *DepositService) Run() {
	service.is_alive.Store(true)
	service.waitgroup.Add(1)
	go service.async_run()
}

func (service *DepositService) Shutdown() {
	service.is_alive.Store(false)
	service.waitgroup.Wait()
}
