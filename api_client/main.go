package main

import (
	"api_client/config"
	"api_client/implementation"
	"log"
	"os"
)

const (
	content_type_json string = "application/json"
	http_timeout      int    = 30 // s
)

func main() {

	if len(os.Args) == 1 {
		implementation.PrintHelpMenu()
		return
	}

	// Load configuration file
	config_file_path := "config.yml"
	config, err := config.Load(config_file_path)
	if err != nil {
		log.Fatal("Unable to load configuration file at ", config_file_path)
	}

	// Run user command. The client will wait for a response from the server before shutting down.
	api_client := implementation.CreateAPIClient(config)
	api_client.Run()
}
