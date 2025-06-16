package implementation

import (
	"api_gateway/config"
	"context"
	"encoding/json"
	"log"
	"shared/messages"
	"shared/responses"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

// Test service for testing the API gateway
type test_service struct {
	config             *config.Service
	is_alive           atomic.Bool
	waitgroup          sync.WaitGroup
	background_context context.Context
	requests_queue     *redis.Client
	responses_queue    *redis.Client
}

func create_test_service(config *config.Service) *test_service {
	service := &test_service{
		config:             config,
		background_context: context.Background(),
	}
	return service
}

func (service *test_service) prepare_redis_clients() error {

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

func (service *test_service) send_response(response_message *responses.Balance) {
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

func (service *test_service) send_failed_response(message string, request_message *messages.GET_Balance) {
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

func (service *test_service) async_run() {
	defer func() {
		service.waitgroup.Done()
		log.Println("Shutdown test service.")
	}()

	log.Println("Started up test service.")

	err := service.prepare_redis_clients()
	if err != nil {
		log.Fatal("Could not connect to Redis server: ", err)
	}
	log.Println("Created clients for Redis message queues.")

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

		// Prepare response
		response_message := responses.Balance{
			Header: responses.Header{
				MessageID: request_message.Header.MessageID,
				Action:    request_message.Header.Action,
			},
			Status: responses.Status_successful,
		}
		switch request_message.Header.Action {
		case messages.Action_get_balance:
			response_message.ErrorMessage = "Test service: I received your message: " + request_message.WalletID
		default:
			service.send_failed_response("Test service: Unknown message received.", &request_message)
			continue
		}
		service.send_response(&response_message)
	}
}

func (service *test_service) run() {
	service.is_alive.Store(true)
	service.waitgroup.Add(1)
	go service.async_run()
}

func (service *test_service) shutdown() {
	service.is_alive.Store(false)
	service.waitgroup.Wait()
}
