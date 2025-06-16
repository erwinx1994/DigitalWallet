package implementation

import (
	"api_gateway/config"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"shared/responses"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type APIGateway struct {
	config           *config.Config
	is_alive         atomic.Bool
	http_server      *http.Server
	http_multiplexer *http_request_multiplexer
	redis_manager    *redis_manager
	waitgroup        sync.WaitGroup
}

func CreateAPIGateway(config *config.Config) (*APIGateway, error) {

	background_context := context.Background()

	redis_manager_, err := create_redis_manager(config, background_context)
	if err != nil {
		return nil, err
	}

	http_multiplexer, err := create_http_request_multiplexer(config, redis_manager_, background_context)
	if err != nil {
		return nil, err
	}

	api_gateway := &APIGateway{
		config:           config,
		http_server:      nil,
		http_multiplexer: http_multiplexer,
		redis_manager:    redis_manager_,
	}

	return api_gateway, nil
}

func (api_gateway *APIGateway) Run() {
	api_gateway.is_alive.Store(true)

	api_gateway.waitgroup.Add(1)
	go api_gateway.async_read_responses(
		service_balance,
		api_gateway.redis_manager.balance_responses_queue,
		&api_gateway.http_multiplexer.balance_responses_cache,
		&api_gateway.config.BalanceService.ResponsesQueue)

	api_gateway.waitgroup.Add(1)
	go api_gateway.async_read_responses(
		service_deposit,
		api_gateway.redis_manager.deposit_responses_queue,
		&api_gateway.http_multiplexer.deposit_responses_cache,
		&api_gateway.config.DepositsService.ResponsesQueue)

	api_gateway.waitgroup.Add(1)
	go api_gateway.async_read_responses(
		service_transaction_history,
		api_gateway.redis_manager.transaction_history_responses_queue,
		&api_gateway.http_multiplexer.transaction_history_responses_cache,
		&api_gateway.config.TransactionHistoryService.ResponsesQueue)

	api_gateway.waitgroup.Add(1)
	go api_gateway.async_read_responses(
		service_transfer,
		api_gateway.redis_manager.transfer_responses_queue,
		&api_gateway.http_multiplexer.transfer_responses_cache,
		&api_gateway.config.TransferService.ResponsesQueue)

	api_gateway.waitgroup.Add(1)
	go api_gateway.async_read_responses(
		service_withdraw,
		api_gateway.redis_manager.withdrawal_responses_queue,
		&api_gateway.http_multiplexer.withdrawal_responses_cache,
		&api_gateway.config.WithdrawalService.ResponsesQueue)

	api_gateway.waitgroup.Add(1)
	go api_gateway.async_http_server()
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

func (api_gateway *APIGateway) async_http_server() {
	defer func() {
		api_gateway.waitgroup.Done()
		log.Println("Shutdown HTTP server.")
	}()

	config := api_gateway.config

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
			Handler:      api_gateway.http_multiplexer,
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

func (api_gateway *APIGateway) async_read_responses(
	service_type int,
	responses_queue *redis.Client,
	response_cache *sync.Map,
	config *config.RedisMessageQueue) {

	service_name := get_service_name(service_type)
	defer func() {
		api_gateway.waitgroup.Done()
		log.Println("Shutdown " + service_name + " responses thread.")
	}()

	log.Println("Started up " + service_name + " responses thread.")

	// Define aliases
	timeout := time.Duration(config.Timeout) * time.Second
	queue_name := config.QueueName

	for api_gateway.is_alive.Load() {

		// Read from Redis responses queue
		timeout_context, cancel := context.WithTimeout(api_gateway.http_multiplexer.context, timeout)
		string_slice, err := responses_queue.BRPop(timeout_context, timeout, queue_name).Result()
		if err != nil {
			cancel()
			continue
		}
		cancel()

		// string_slice[0] gives the name of the queue
		// string_slice[1] gives the data retrieved from the queue
		// Response is prepared by the primary backend service, not the API gateway.

		// Put the response in the asynchronous map
		switch service_type {
		case service_balance:

			response_message := responses.Balance{}
			bytes := []byte(string_slice[1])
			err := json.Unmarshal(bytes, &response_message)
			if err != nil {
				log.Println("Error deserialising JSON message. Must not happen in production.")
				continue
			}
			response_cache.Store(response_message.Header.MessageID, bytes)

		case service_deposit:

			response_message := responses.Deposit{}
			bytes := []byte(string_slice[1])
			err := json.Unmarshal(bytes, &response_message)
			if err != nil {
				log.Println("Error deserialising JSON message. Must not happen in production.")
				continue
			}
			response_cache.Store(response_message.Header.MessageID, bytes)

		case service_transaction_history:

			response_message := responses.TransactionHistory{}
			bytes := []byte(string_slice[1])
			err := json.Unmarshal(bytes, &response_message)
			if err != nil {
				log.Println("Error deserialising JSON message. Must not happen in production.")
				continue
			}
			response_cache.Store(response_message.Header.MessageID, bytes)

		case service_transfer:

			response_message := responses.Transfer{}
			bytes := []byte(string_slice[1])
			err := json.Unmarshal(bytes, &response_message)
			if err != nil {
				log.Println("Error deserialising JSON message. Must not happen in production.")
				continue
			}
			response_cache.Store(response_message.Header.MessageID, bytes)

		case service_withdraw:

			response_message := responses.Withdraw{}
			bytes := []byte(string_slice[1])
			err := json.Unmarshal(bytes, &response_message)
			if err != nil {
				log.Println("Error deserialising JSON message. Must not happen in production.")
				continue
			}
			response_cache.Store(response_message.Header.MessageID, bytes)
		}
	}
}
