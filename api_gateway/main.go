package main

import (
	"api_gateway/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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

	// Start running API gateway
	api_gateway := CreateAPIGateway(config)
	api_gateway.run()

	// Listen for abort signal to terminate the api_gateway
	// Pressing CTRL + C while the application is running
	abort_channel := make(chan os.Signal, 1)
	signal.Notify(abort_channel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-abort_channel

	// Shutdown the HTTP server gracefully
	api_gateway.shutdown()
}
