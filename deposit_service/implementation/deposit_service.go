package implementation

import (
	"context"
	"deposit_service/config"
	"encoding/json"
	"log"
	"shared/messages"
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
		redis_message := messages.RedisMessage{
			Body: messages.POST_Deposit{},
		}
		err = json.Unmarshal([]byte(string_slice[1]), &redis_message)
		if err != nil {
			log.Println("Failed to deserialise JSON message. Should not happen in production.")
			// In practice, we will need an error notification system. Building an error
			// notification is skipped due to time constraints.
			continue
		}

		// Query postgre SQL database

		// Put response into response queue

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
