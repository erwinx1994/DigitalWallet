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
)

type APIGateway struct {
	config           *config.Config
	is_alive         atomic.Bool
	http_server      *http.Server
	http_multiplexer HTTPRequestMultiplexer
	waitgroup        sync.WaitGroup
}

func CreateAPIGateway(config *config.Config) *APIGateway {
	api_gateway := &APIGateway{
		config:      config,
		http_server: nil,
	}
	return api_gateway
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
