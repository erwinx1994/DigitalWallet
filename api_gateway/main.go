package main

import (
	"api_gateway/config"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type HTTPRequestMultiplexer struct{}

func (mux *HTTPRequestMultiplexer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Print("HTTP Request received.")
}

func main() {

	log.Println("DIGITAL WALLET INC API GATEWAY")
	log.Println("This API gateway is a backend end point for the user to interact with digital wallet services provided by Digital Wallet Inc.")

	// Load configuration file
	config_file_path := "config.yml"
	if len(os.Args) > 1 {
		config_file_path = os.Args[1]
	}
	config, err := config.Load(config_file_path)
	if err != nil {
		log.Fatal("Unable to load configuration file at ", config_file_path)
	}
	log.Println("Successfully loaded configuration file at ", config_file_path)

	var is_alive atomic.Bool
	is_alive.Store(true)

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)

	var http_server *http.Server

	go func() {

		defer func() {
			waitgroup.Done()
			log.Println("Shutdown HTTP server.")
		}()

		// Server continues running until terminated by user
		for is_alive.Load() {

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

			// Create a HTTP request multiplexer
			mux := HTTPRequestMultiplexer{}

			http_server = &http.Server{
				Addr:                         ":" + config.ListenPort,
				Handler:                      &mux,
				DisableGeneralOptionsHandler: false,
				ReadTimeout:                  time.Duration(config.ReadTimeout),
				WriteTimeout:                 time.Duration(config.WriteTimeout),
				IdleTimeout:                  time.Duration(config.IdleTimeout),
			}
			err = http_server.Serve(listener)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Print("Error with HTTP server: ", err.Error())
				log.Print("Reestablishing in ", config.RetryInterval, " s.")
				time.Sleep(time.Duration(config.RetryInterval) * time.Second)
				continue
			}
		}

	}()

	// Listen for abort signal to terminate the api_gateway
	// Pressing CTRL + C while the application is running
	abort_channel := make(chan os.Signal, 1)
	signal.Notify(abort_channel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-abort_channel

	// Shutdown the HTTP server gracefully
	is_alive.Store(false)
	if http_server != nil {
		// No need to handle error returned by Server.Shutdown.
		// The signal to abort is already sent. Just terminate the application
		// regardless of whether an error occurs during shutdown.
		context_ := context.Background()
		context_with_timeout, cancel := context.WithTimeout(context_, time.Duration(5)*time.Second)
		defer cancel()
		http_server.Shutdown(context_with_timeout)
	}
	waitgroup.Wait()
}
