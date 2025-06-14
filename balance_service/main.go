package main

import (
	"balance_service/config"
	"balance_service/implementation"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Println("DIGITAL WALLET INC BALANCE SERVICE")
	log.Println("The balance service retrieves the balance of the user's wallet and returns it as a response.")

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

	// Start running balance service
	balance_service := implementation.CreateBalanceService(config)
	balance_service.Run()

	// Listen for abort signal to terminate the balance service
	// Pressing CTRL + C while the application is running
	abort_channel := make(chan os.Signal, 1)
	signal.Notify(abort_channel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-abort_channel

	// Shutdown the HTTP server gracefully
	balance_service.Shutdown()
}
