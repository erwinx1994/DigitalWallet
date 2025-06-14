package implementation

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"transfer_service/config"

	"github.com/redis/go-redis/v9"
)

type TransferService struct {
	config          *config.Config
	is_alive        atomic.Bool
	waitgroup       sync.WaitGroup
	requests_queue  *redis.Client
	responses_queue *redis.Client
}

func CreateTransferService(config *config.Config) *TransferService {
	service := &TransferService{
		config: config,
	}
	return service
}

func (service *TransferService) prepare_redis_clients() error {

	background_context := context.Background()

	{
		// Prepare requests queue
		requests_queue := redis.NewClient(service.config.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(service.config.RequestsQueue.Timeout)*time.Second)
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

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(service.config.ResponsesQueue.Timeout)*time.Second)
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

func (service *TransferService) async_run() {
	defer func() {
		service.waitgroup.Done()
		log.Println("Shutdown transfer service.")
	}()

	err := service.prepare_redis_clients()
	if err != nil {
		log.Fatal("Could not connect to Redis server: ", err)
	}

	// Service continues running until terminated by user
	for service.is_alive.Load() {

		log.Println("Started up transfer service.")
		time.Sleep(1 * time.Second)
	}
}

func (service *TransferService) Run() {
	service.is_alive.Store(true)
	service.waitgroup.Add(1)
	go service.async_run()
}

func (service *TransferService) Shutdown() {
	service.is_alive.Store(false)
	service.waitgroup.Wait()
}
