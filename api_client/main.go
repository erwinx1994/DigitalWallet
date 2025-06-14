package main

import (
	"api_client/config"
	"api_client/implementation"
	"fmt"
	"log"
	"os"
)

const (
	content_type_json string = "application/json"
	http_timeout      int    = 30 // s
)

func print_help_menu() {

	fmt.Println("DIGITAL WALLET INC CLIENT")
	fmt.Println()

	fmt.Println("This command allows you to deposit an amount of money in the specified wallet in the specified currency.")
	fmt.Println()

	fmt.Println("\tapi_client deposit <wallet_id> <currency> <amount>")
	fmt.Println()

	fmt.Println("This command allows you to withdraw an amount of money from the specified wallet in the specified currency.")
	fmt.Println()

	fmt.Println("\tapi_client withdraw <wallet_id> <currency> <amount>")
	fmt.Println()

	fmt.Println("This command allows you to transfer an amount of money from source to destination wallets of the same currency.")
	fmt.Println()

	fmt.Println("\tapi_client transfer <source_wallet_id> <destination_wallet_id> <currency> <amount>")
	fmt.Println()

	fmt.Println("This command allows you to get the balance of the specified wallet.")
	fmt.Println()

	fmt.Println("\tapi_client get_balance <wallet_id>")
	fmt.Println()

	fmt.Println("This command allows you to get the transaction history of the specified wallet from a start date to end date, inclusive.")
	fmt.Println()

	fmt.Println("\tapi_client get_transaction_history <wallet_id> <start_date> <end_date>")
	fmt.Println()
}

func main() {

	if len(os.Args) == 1 {
		print_help_menu()
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
