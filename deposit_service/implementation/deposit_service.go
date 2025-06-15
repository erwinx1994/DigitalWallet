package implementation

import (
	"context"
	"database/sql"
	"deposit_service/config"
	"encoding/json"
	"log"
	"shared/messages"
	"shared/responses"
	"sync"
	"sync/atomic"
	"time"

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

		// Query PostgreSQL database

		// Prepare response
		response_message := responses.Deposit{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
		}
		bytes_to_send, err := json.Marshal(response_message)
		if err != nil {
			log.Println("Failed to serialise response message. Should not happen in production.")
			// In practice, we will need an error notification system. I have skipped
			// building an error notification system due to time constraints.
			continue
		}

		// Put response into responses queue
		timeout = time.Duration(service.config.ResponsesQueue.Timeout) * time.Second
		queue_name = service.config.ResponsesQueue.QueueName
		timeout_context, cancel = context.WithTimeout(service.background_context, timeout)
		_, err = service.responses_queue.LPush(timeout_context, queue_name, bytes_to_send).Result()
		if err != nil {
			cancel()
			log.Println("Failed to put response into responses queue.")
			return
		}
		cancel()

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
