package implementation

import (
	"balance_service/config"
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

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type BalanceService struct {
	config             *config.Config
	is_alive           atomic.Bool
	waitgroup          sync.WaitGroup
	background_context context.Context
	requests_queue     *redis.Client
	responses_queue    *redis.Client
}

func CreateBalanceService(config *config.Config) *BalanceService {
	service := &BalanceService{
		config:             config,
		background_context: context.Background(),
	}
	return service
}

func (service *BalanceService) prepare_redis_clients() error {

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

func (service *BalanceService) send_response(response_message *responses.Balance) {
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

func (service *BalanceService) send_failed_response(message string, request_message *messages.GET_Balance) {
	response_message := responses.Balance{
		Header: responses.Header{
			MessageID: request_message.Header.MessageID,
			Action:    request_message.Header.Action,
		},
		Status:       responses.Status_failed,
		ErrorMessage: message,
	}
	service.send_response(&response_message)
}

func (service *BalanceService) async_run() {
	defer func() {
		service.waitgroup.Done()
		log.Println("Shutdown balance service.")
	}()

	log.Println("Started up balance service.")

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
	get_balance, err := db.Prepare("select currency, balance from " + service.config.WalletDatabase.BalanceTable + " where wallet_id=$1")
	if err != nil {
		log.Fatal("Unable to prepare SQL statement.")
	}
	defer get_balance.Close()

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
		request_message := messages.GET_Balance{}
		err = json.Unmarshal([]byte(string_slice[1]), &request_message)
		if err != nil {
			log.Println("Failed to deserialise JSON message. Should not happen in production.")
			// In practice, we will need an error notification system. I have skipped
			// building an error notification system due to time constraints.
			continue
		}

		// Verify that the correct message was received
		if request_message.Header.Action != messages.Action_get_balance {
			service.send_failed_response("Message received by wrong service", &request_message)
			continue
		}

		// Query PostgreSQL database. Assume inputs are correct.
		var currency string = ""
		var balance int64 = 0
		err = get_balance.QueryRow(request_message.WalletID).Scan(&currency, &balance)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				service.send_failed_response("Wallet does not exist", &request_message)
			} else {
				service.send_failed_response("Database error", &request_message)
			}
			continue
		}

		// Prepare response
		response_message := responses.Balance{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
			Status:   responses.Status_successful,
			Currency: currency,
			Balance:  utilities.Convert_database_to_display_format(balance),
		}
		service.send_response(&response_message)
	}
}

func (service *BalanceService) Run() {
	service.is_alive.Store(true)
	service.waitgroup.Add(1)
	go service.async_run()
}

func (service *BalanceService) Shutdown() {
	service.is_alive.Store(false)
	service.waitgroup.Wait()
}
