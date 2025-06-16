package implementation

import (
	"api_client/config"
	"fmt"
	"net/http"
	"os"
	"time"
)

type APIClient struct {
	config *config.Config
}

func CreateAPIClient(config *config.Config) *APIClient {
	api_client := &APIClient{
		config: config,
	}
	return api_client
}

func (api_client *APIClient) post_deposit() {

	// Verify that inputs are correct
	if len(os.Args) != number_of_arguments_deposit {
		fmt.Println("Incorrect number of arguments for deposit command.")
		return
	}

	// Send request to server

	// Wait for response from server

	// Print result to console

}

func (api_client *APIClient) post_withdrawal() {

	// Verify that inputs are correct
	if len(os.Args) != number_of_arguments_withdraw {
		fmt.Println("Incorrect number of arguments for withdraw command.")
		return
	}

	// Send request to server

	// Wait for response from server

	// Print result to console

}

func (api_client *APIClient) post_transfer() {

	// Verify that inputs are correct
	if len(os.Args) != number_of_arguments_transfer {
		fmt.Println("Incorrect number of arguments for transfer command.")
		return
	}

	// Send request to server

	// Wait for response from server

	// Print result to console

}

func (api_client *APIClient) get_wallet_balance() {

	// Verify that inputs are correct
	if len(os.Args) != number_of_arguments_get_balance {
		fmt.Println("Incorrect number of arguments for get_balance command. Please review the help menu for assistance. It can be accessed just by entering api_client.")
		return
	}
	if len(os.Args[2]) == 0 {
		fmt.Println("Please enter a wallet ID.")
		fmt.Println()
		fmt.Println("api_client get_balance <wallet_id>")
		return
	}

	// api_client get_balance <wallet_id>
	// GET /wallets/{wallet_id}/balance

	// Prepare GET request
	http_client := http.Client{
		Timeout: time.Duration(api_client.config.RequestTimeout) * time.Second,
	}
	wallet_id := os.Args[2]
	base_url := api_client.config.Server.GetURL()
	full_url := base_url + "/wallets/{" + wallet_id + "}/balance"
	response, err := http_client.Get(full_url)
	if err != nil {
		fmt.Println("HTTP error occurred: ", err.Error())
		return
	}

	// Parse result
	bytes := make([]byte, response.ContentLength)
	_, err = response.Body.Read(bytes)
	if err != nil {
		fmt.Println("Error reading response: ", err.Error())
		return
	}

	// Print result to console
	fmt.Println("HTTP GET response status: ", response.Status)
	fmt.Println(string(bytes))
}

func (api_client *APIClient) get_transaction_history() {

	// Verify that inputs are correct
	if len(os.Args) != number_of_arguments_get_transaction_history {
		fmt.Println("Incorrect number of arguments for get_transaction_history command.")
		return
	}

	// Send request to server

	// Wait for response from server

	// Print result to console

}

func (api_client *APIClient) Run() {
	verb := os.Args[1]
	switch verb {
	case action_deposit:
		api_client.post_deposit()
	case action_withdraw:
		api_client.post_withdrawal()
	case action_transfer:
		api_client.post_transfer()
	case action_get_balance:
		api_client.get_wallet_balance()
	case action_get_transaction_history:
		api_client.get_transaction_history()
	default:
		fmt.Println("Invalid command. Please review the help menu for assistance. It can be accessed by entering this command without any arguments.")
		fmt.Println()
		fmt.Println("api_client")
	}
}
