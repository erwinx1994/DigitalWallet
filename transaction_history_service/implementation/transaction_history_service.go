package implementation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"shared/messages"
	"shared/responses"
	"shared/utilities"
	"sync"
	"sync/atomic"
	"time"
	"transaction_history_service/config"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	time_format string = "20060102"
)

type TransactionHistoryService struct {
	config             *config.Config
	is_alive           atomic.Bool
	waitgroup          sync.WaitGroup
	background_context context.Context
	requests_queue     *redis.Client
	responses_queue    *redis.Client
}

func CreateTransactionHistoryService(config *config.Config) *TransactionHistoryService {
	service := &TransactionHistoryService{
		config:             config,
		background_context: context.Background(),
	}
	return service
}

func (service *TransactionHistoryService) prepare_redis_clients() error {

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

func (service *TransactionHistoryService) send_response(response_message *responses.TransactionHistory) {
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

func (service *TransactionHistoryService) send_failed_response(message string, request_message *messages.GET_TransactionHistory) {
	response_message := responses.TransactionHistory{
		Header: responses.Header{
			MessageID: request_message.Header.MessageID,
			Action:    request_message.Header.Action,
		},
		Status:       responses.Status_failed,
		ErrorMessage: message,
	}
	service.send_response(&response_message)
}

func (service *TransactionHistoryService) async_run() {
	defer func() {
		service.waitgroup.Done()
		log.Println("Shutdown transaction history service.")
	}()

	log.Println("Started up transaction history service.")

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
	get_balance, err := db.Prepare("select balance from " + service.config.WalletDatabase.BalanceTable + " where wallet_id=$1")
	if err != nil {
		log.Fatal("Unable to prepare SQL statement.")
	}
	defer get_balance.Close()

	get_transaction_history, err := db.Prepare("select date_and_time, currency, amount from " + service.config.WalletDatabase.TransactionsTable + " where wallet_id=$1 and date_and_time>=$2 and date_and_time<$3 order by date_and_time desc")
	if err != nil {
		log.Fatal("Unable to prepare SQL statement.")
	}
	defer get_transaction_history.Close()

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
		request_message := messages.GET_TransactionHistory{}
		err = json.Unmarshal([]byte(string_slice[1]), &request_message)
		if err != nil {
			log.Println("Failed to deserialise JSON message. Should not happen in production.")
			// In practice, we will need an error notification system. Building an error
			// notification is skipped due to time constraints.
			continue
		}

		// Verify that the correct message was received
		if request_message.Header.Action != messages.Action_get_transaction_history {
			service.send_failed_response("Message received by wrong service", &request_message)
			continue
		}

		// Query PostgreSQL database
		db_transaction, err := db.Begin()
		if err != nil {
			service.send_failed_response("Database error", &request_message)
			continue
		}

		// Check if wallet already exist, return an error if it does not
		var balance int64 = 0
		tx_get_balance := db_transaction.Stmt(get_balance)
		err = tx_get_balance.QueryRow(request_message.WalletID).Scan(&balance)
		if err != nil {
			db_transaction.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				service.send_failed_response("Cannot get transaction history of non-existent wallet", &request_message)
			} else {
				service.send_failed_response("Database error", &request_message)
			}
			continue
		}

		// Formulate parameters
		from := time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC)
		to := time.Now().UTC()
		if len(request_message.From) > 0 {
			from, err = time.Parse(time_format, request_message.From)
			if err != nil {
				db_transaction.Rollback()
				service.send_failed_response("Invalid start date", &request_message)
				continue
			}
		}
		if len(request_message.To) > 0 {
			to, err = time.Parse(time_format, request_message.To)
			if err != nil {
				db_transaction.Rollback()
				service.send_failed_response("Invalid end date", &request_message)
				continue
			}
			to = to.AddDate(0, 0, 1)
		}

		// Get transaction history of wallet
		tx_get_transaction_history := db_transaction.Stmt(get_transaction_history)
		rows, err := tx_get_transaction_history.Query(request_message.WalletID, from, to)
		if err != nil {
			db_transaction.Rollback()
			service.send_failed_response("Database error", &request_message)
			continue
		}
		transaction_history := []responses.Transaction{}
		var date_and_time time.Time
		var currency string = ""
		var amount int64 = 0
		for rows.Next() {
			err := rows.Scan(&date_and_time, &currency, &amount)
			if err != nil {
				service.send_failed_response("Database error", &request_message)
				break
			}

			transaction_type := responses.Transaction_type_unknown
			if amount >= 0 {
				transaction_type = responses.Transaction_type_deposit
			} else if amount < 0 {
				transaction_type = responses.Transaction_type_withdraw
			}

			if amount < 0 {
				amount *= -1
			}

			transaction_history = append(
				transaction_history,
				responses.Transaction{
					Date:     date_and_time.Format(time_format),
					Type:     transaction_type,
					Currency: currency,
					Amount:   utilities.Convert_database_to_display_format(amount),
				})
		}
		err = rows.Err()
		if err != nil {
			rows.Close()
			db_transaction.Rollback()
			service.send_failed_response("Database error", &request_message)
			continue
		}
		rows.Close()
		db_transaction.Commit()

		// Prepare response
		response_message := responses.TransactionHistory{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
			Status:  responses.Status_successful,
			History: transaction_history,
		}
		service.send_response(&response_message)
	}
}

func (service *TransactionHistoryService) Run() {
	service.is_alive.Store(true)
	service.waitgroup.Add(1)
	go service.async_run()
}

func (service *TransactionHistoryService) Shutdown() {
	service.is_alive.Store(false)
	service.waitgroup.Wait()
}
