package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"withdraw_service/config"
	"withdraw_service/implementation"
)

func main() {

	log.Println("DIGITAL WALLET INC WITHDRAW SERVICE")
	log.Println("The withdraw service creates a new withdrawal transaction for a wallet and updates its balance.")

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

	// Start running withdrawal service
	service := implementation.CreateWithdrawService(config)
	service.Run()

	// Listen for abort signal to terminate the balance service
	// Pressing CTRL + C while the application is running
	abort_channel := make(chan os.Signal, 1)
	signal.Notify(abort_channel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-abort_channel

	// Shutdown the withdrawal service gracefully
	service.Shutdown()
}
