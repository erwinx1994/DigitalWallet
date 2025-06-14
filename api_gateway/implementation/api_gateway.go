package implementation

import (
	"api_gateway/config"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type APIGateway struct {

	// API gateway settings
	config           *config.Config
	is_alive         atomic.Bool
	http_server      *http.Server
	http_multiplexer HTTPRequestMultiplexer
	waitgroup        sync.WaitGroup

	// Resource handles to request and response queues
	deposit_requests_queue              *redis.Client
	deposit_responses_queue             *redis.Client
	withdrawal_requests_queue           *redis.Client
	withdrawal_responses_queue          *redis.Client
	transfer_requests_queue             *redis.Client
	transfer_responses_queue            *redis.Client
	balance_requests_queue              *redis.Client
	balance_responses_queue             *redis.Client
	transaction_history_requests_queue  *redis.Client
	transaction_history_responses_queue *redis.Client
}

func CreateAPIGateway(config *config.Config) *APIGateway {
	api_gateway := &APIGateway{
		config:      config,
		http_server: nil,
	}
	return api_gateway
}

func (api_gateway *APIGateway) prepare_redis_clients() error {

	background_context := context.Background()

	{
		// Prepare deposit requests queue
		api_gateway.deposit_requests_queue = redis.NewClient(api_gateway.config.DepositsService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.DepositsService.RequestsQueue.Timeout)*time.Second)
		_, err := api_gateway.deposit_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare deposit responses queue
		api_gateway.deposit_responses_queue = redis.NewClient(api_gateway.config.DepositsService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.DepositsService.ResponsesQueue.Timeout)*time.Second)
		_, err := api_gateway.deposit_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare withdrawal requests queue
		api_gateway.withdrawal_requests_queue = redis.NewClient(api_gateway.config.WithdrawalService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.WithdrawalService.RequestsQueue.Timeout)*time.Second)
		_, err := api_gateway.withdrawal_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare withdrawal responses queue
		api_gateway.withdrawal_responses_queue = redis.NewClient(api_gateway.config.WithdrawalService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.WithdrawalService.ResponsesQueue.Timeout)*time.Second)
		_, err := api_gateway.withdrawal_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare transfer requests queue
		api_gateway.transfer_requests_queue = redis.NewClient(api_gateway.config.TransferService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.TransferService.RequestsQueue.Timeout)*time.Second)
		_, err := api_gateway.transfer_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare transfer responses queue
		api_gateway.transfer_responses_queue = redis.NewClient(api_gateway.config.TransferService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.TransferService.ResponsesQueue.Timeout)*time.Second)
		_, err := api_gateway.transfer_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare balance requests queue
		api_gateway.balance_requests_queue = redis.NewClient(api_gateway.config.BalanceService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.BalanceService.RequestsQueue.Timeout)*time.Second)
		_, err := api_gateway.balance_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare balance responses queue
		api_gateway.balance_responses_queue = redis.NewClient(api_gateway.config.BalanceService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.BalanceService.ResponsesQueue.Timeout)*time.Second)
		_, err := api_gateway.balance_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare transaction history requests queue
		api_gateway.transaction_history_requests_queue = redis.NewClient(api_gateway.config.TransactionHistoryService.RequestsQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.TransactionHistoryService.RequestsQueue.Timeout)*time.Second)
		_, err := api_gateway.transaction_history_requests_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	{
		// Prepare transaction history responses queue
		api_gateway.transaction_history_responses_queue = redis.NewClient(api_gateway.config.TransactionHistoryService.ResponsesQueue.GetRedisOptions())

		timeout_context, cancel := context.WithTimeout(background_context, time.Duration(api_gateway.config.TransactionHistoryService.ResponsesQueue.Timeout)*time.Second)
		_, err := api_gateway.transaction_history_responses_queue.Ping(timeout_context).Result()
		if err != nil {
			cancel()
			return err
		}
		cancel()
	}

	return nil
}

func (api_gateway *APIGateway) async_http_server() {
	defer func() {
		api_gateway.waitgroup.Done()
		log.Println("Shutdown HTTP server.")
	}()

	config := api_gateway.config

	err := api_gateway.prepare_redis_clients()
	if err != nil {
		log.Fatal("Could not connect to Redis server: ", err)
	}

	// Server continues running until terminated by user
	for api_gateway.is_alive.Load() {

		log.Println("Starting up HTTP server.")

		// Create a listener
		listener, err := net.Listen("tcp", ":"+config.HTTPServer.ListenPort)
		if err != nil {
			log.Print("Unable to listen on port: ", config.HTTPServer.ListenPort)
			log.Print("Reestablishing in ", config.HTTPServer.RetryInterval, " s.")
			time.Sleep(time.Duration(config.HTTPServer.RetryInterval) * time.Second)
			continue
		}
		log.Println("Listening on port: ", config.HTTPServer.ListenPort)

		api_gateway.http_server = &http.Server{
			Addr:         ":" + config.HTTPServer.ListenPort,
			Handler:      &api_gateway.http_multiplexer,
			ReadTimeout:  time.Duration(config.HTTPServer.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(config.HTTPServer.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(config.HTTPServer.IdleTimeout) * time.Second,
		}
		err = api_gateway.http_server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Print("Error with HTTP server: ", err.Error())
			log.Print("Reestablishing in ", config.HTTPServer.RetryInterval, " s.")
			time.Sleep(time.Duration(config.HTTPServer.RetryInterval) * time.Second)
			continue
		}
	}
}

func (api_gateway *APIGateway) async_read_responses() {

	defer func() {
		api_gateway.waitgroup.Done()
		log.Println("Shutdown API responses thread.")
	}()

	log.Println("Started up response reading thread.")
	for api_gateway.is_alive.Load() {

		// Read from Redis message queues

		// Put the responses in api_gateway.http_multiplexer.api_responses

		time.Sleep(1 * time.Second)
	}
}

func (api_gateway *APIGateway) Run() {
	api_gateway.is_alive.Store(true)
	api_gateway.waitgroup.Add(2)
	go api_gateway.async_http_server()
	go api_gateway.async_read_responses()
}

func (api_gateway *APIGateway) Shutdown() {
	api_gateway.is_alive.Store(false)
	if api_gateway.http_server != nil {
		// No need to handle error returned by Server.Shutdown.
		// The signal to abort is already sent. Just terminate the application
		// regardless of whether an error occurs during shutdown.
		context_ := context.Background()
		context_with_timeout, cancel := context.WithTimeout(context_, time.Duration(5)*time.Second)
		defer cancel()
		api_gateway.http_server.Shutdown(context_with_timeout)
	}
	api_gateway.waitgroup.Wait()
}
