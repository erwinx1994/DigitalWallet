package main

import (
	"api_gateway/config"
	"api_gateway/paths"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type QueueItem struct {
	id        int64
	json_body []byte
}

// A custom HTTP request multiplexer is needed as the multiplexer also needs to store the responses
// to each request in a map which will be populated by another thread.
type HTTPRequestMultiplexer struct {
	api_responses sync.Map
}

func (mux *HTTPRequestMultiplexer) POST_Deposit(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) POST_Withdrawal(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) POST_Transfer(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) GET_WalletBalance(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) GET_TransactionHistory(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) GET_Test(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Log content of request
	log.Println("Method: ", request.Method)
	log.Println("Scheme: ", request.URL.Scheme)
	log.Println("Host: ", request.URL.Host)
	log.Println("Path: ", request.URL.Path)
	log.Println("Header: ", request.Header)

	// Send test response to user
	type TestResponse struct {
		Message string `json:",omitempty"`
	}
	test_response := TestResponse{
		Message: "Hi. This is a test response.",
	}
	bytes, err := json.Marshal(test_response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write([]byte(bytes))
}

func (mux *HTTPRequestMultiplexer) ProcessGETRequests(writer http.ResponseWriter, request *http.Request) {

	// Request for balance
	result := paths.MatchAndExtract(request.URL.Path, paths.Wallets_balance)
	if result.MatchFound {
		mux.GET_WalletBalance(result, writer, request)
		return
	}

	// Request for transaction history
	result = paths.MatchAndExtract(request.URL.Path, paths.Wallets_transaction_history)
	if result.MatchFound {
		mux.GET_TransactionHistory(result, writer, request)
		return
	}

	// Test request
	result = paths.MatchAndExtract(request.URL.Path, paths.Test)
	if result.MatchFound {
		mux.GET_Test(result, writer, request)
		return
	}
}

func (mux *HTTPRequestMultiplexer) ProcessPOSTRequests(writer http.ResponseWriter, request *http.Request) {

	// Deposit money into wallet
	result := paths.MatchAndExtract(request.URL.Path, paths.Wallets_deposits)
	if result.MatchFound {
		mux.POST_Deposit(result, writer, request)
		return
	}

	// Withdraw money from wallet
	result = paths.MatchAndExtract(request.URL.Path, paths.Wallets_withdrawals)
	if result.MatchFound {
		mux.POST_Withdrawal(result, writer, request)
		return
	}

	// Transfer money from one wallet to another
	result = paths.MatchAndExtract(request.URL.Path, paths.Transfer)
	if result.MatchFound {
		mux.POST_Transfer(result, writer, request)
		return
	}

}

// A new goroutine is created to serve each HTTP request
func (mux *HTTPRequestMultiplexer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		mux.ProcessGETRequests(writer, request)
	} else if request.Method == http.MethodPost {
		mux.ProcessPOSTRequests(writer, request)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type APIGateway struct {
	config           *config.Config
	is_alive         atomic.Bool
	http_server      *http.Server
	http_multiplexer HTTPRequestMultiplexer
	waitgroup        sync.WaitGroup
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
		listener, err := net.Listen("tcp", ":"+config.ListenPort)
		if err != nil {
			log.Print("Unable to listen on port: ", config.ListenPort)
			log.Print("Reestablishing in ", config.RetryInterval, " s.")
			time.Sleep(time.Duration(config.RetryInterval) * time.Second)
			continue
		}
		log.Println("Listening on port: ", config.ListenPort)

		api_gateway.http_server = &http.Server{
			Addr:         ":" + config.ListenPort,
			Handler:      &api_gateway.http_multiplexer,
			ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
		}
		err = api_gateway.http_server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Print("Error with HTTP server: ", err.Error())
			log.Print("Reestablishing in ", config.RetryInterval, " s.")
			time.Sleep(time.Duration(config.RetryInterval) * time.Second)
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

func (api_gateway *APIGateway) run() {
	api_gateway.is_alive.Store(true)
	api_gateway.waitgroup.Add(2)
	go api_gateway.async_http_server()
	go api_gateway.async_read_responses()
}

func (api_gateway *APIGateway) shutdown() {
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
