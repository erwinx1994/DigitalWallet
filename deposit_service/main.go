package main

import (
	"deposit_service/config"
	"deposit_service/implementation"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Println("DIGITAL WALLET INC DEPOSIT SERVICE")
	log.Println("The deposit service creates a new deposit transaction for a wallet and updates its balance.")

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

	// Start running deposit service
	service := implementation.CreateDepositService(config)
	service.Run()

	// Listen for abort signal to terminate the balance service
	// Pressing CTRL + C while the application is running
	abort_channel := make(chan os.Signal, 1)
	signal.Notify(abort_channel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-abort_channel

	// Shutdown the deposit service gracefully
	service.Shutdown()
}
