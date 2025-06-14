package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"transfer_service/config"
	"transfer_service/implementation"
)

func main() {

	log.Println("DIGITAL WALLET INC TRANSFER SERVICE")
	log.Println("The transfer service transfers money from one wallet to another.")

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

	// Start running transfer service
	service := implementation.CreateTransferService(config)
	service.Run()

	// Listen for abort signal to terminate the balance service
	// Pressing CTRL + C while the application is running
	abort_channel := make(chan os.Signal, 1)
	signal.Notify(abort_channel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-abort_channel

	// Shutdown the transfer service gracefully
	service.Shutdown()
}
