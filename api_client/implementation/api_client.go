package implementation

import (
	"api_client/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"shared/messages"
	"shared/responses"
	"time"
)

func ConvertToString(status int) string {
	result := ""
	switch status {
	case responses.Status_unknown:
		result = "Unknown"
	case responses.Status_successful:
		result = "Successful"
	case responses.Status_failed:
		result = "Failed"
	}
	return result
}

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

	/*
		api_client deposit <wallet_id> <currency> <amount>
		POST /wallets/{wallet_id}/deposits
		{
		    "amount": 5000,
		    "currency": "XXX"
		}
	*/

	// Verify that inputs are correct
	if len(os.Args) != number_of_arguments_deposit {
		fmt.Println("Incorrect number of arguments for deposit command. Please review the help menu for assistance. It can be accessed just by entering api_client.")
		return
	}
	if len(os.Args[2]) == 0 {
		fmt.Println("Please enter a wallet ID.")
		fmt.Println()
		fmt.Println("api_client deposit <wallet_id> <currency> <amount>")
		return
	}
	if len(os.Args[3]) == 0 {
		fmt.Println("Please enter currency to deposit.")
		fmt.Println()
		fmt.Println("api_client deposit <wallet_id> <currency> <amount>")
		return
	}
	if len(os.Args[4]) == 0 {
		fmt.Println("Please enter an amount to deposit.")
		fmt.Println()
		fmt.Println("api_client deposit <wallet_id> <currency> <amount>")
		return
	}

	// Make POST request
	http_client := http.Client{
		Timeout: time.Duration(api_client.config.RequestTimeout) * time.Second,
	}
	wallet_id := os.Args[2]
	base_url := api_client.config.Server.GetURL()
	full_url := base_url + "/wallets/{" + wallet_id + "}/deposits"
	http_post_body := messages.POST_Deposit{
		Amount:   os.Args[4],
		Currency: os.Args[3],
	}
	http_post_body_bytes, err := json.Marshal(http_post_body)
	if err != nil {
		fmt.Println("Unable to serialise body of request into JSON.")
		return
	}
	response, err := http_client.Post(full_url, "application/json", bytes.NewReader(http_post_body_bytes))
	if err != nil {
		fmt.Println("HTTP error occurred: ", err.Error())
		return
	}

	// Parse result
	response_bytes := make([]byte, response.ContentLength)
	_, err = response.Body.Read(response_bytes)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			fmt.Println("Error reading response: ", err.Error())
			return
		}
	}
	response_body := responses.Deposit{}
	err = json.Unmarshal(response_bytes, &response_body)
	if err != nil {
		fmt.Println("Error parsing JSON response.")
		return
	}

	// Print result to console
	fmt.Println("Request status: ", ConvertToString(response_body.Status))
	switch response_body.Status {
	case responses.Status_successful:
		fmt.Println("New balance: ", response_body.Currency, " ", response_body.NewBalance)
	case responses.Status_failed:
		fmt.Println("Error message: ", response_body.ErrorMessage)
	case responses.Status_unknown:
	}
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

	// api_client get_balance <wallet_id>
	// GET /wallets/{wallet_id}/balance

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
		if !errors.Is(err, io.EOF) {
			fmt.Println("Error reading response: ", err.Error())
			return
		}
	}
	response_body := responses.Balance{}
	err = json.Unmarshal(bytes, &response_body)
	if err != nil {
		fmt.Println("Error parsing JSON response.")
		return
	}

	// Print result to console
	fmt.Println("Request status: ", ConvertToString(response_body.Status))
	switch response_body.Status {
	case responses.Status_successful:
		fmt.Println("Balance: ", response_body.Currency, " ", response_body.Balance)
	case responses.Status_failed:
		fmt.Println("Error message: ", response_body.ErrorMessage)
	case responses.Status_unknown:
	}

}

func (api_client *APIClient) get_transaction_history() {

	// api_client get_transaction_history <wallet_id> <start_date> <end_date>
	// GET /wallets/{wallet_id}/transaction_history?from=YYYYMMDD&to=YYYYMMDD

	// Verify that inputs are correct
	if len(os.Args) < minimum_number_of_arguments_get_transaction_history {
		fmt.Println("Incorrect number of arguments for get_transaction_history command.")
		return
	}

	if len(os.Args[2]) == 0 {
		fmt.Println("Please enter a wallet ID.")
		fmt.Println()
		fmt.Println("api_client get_transaction_history <wallet_id> <start_date> <end_date>")
		return
	}

	// Could do simple argument checking here. But due to lack of time, error checking
	// is left to the transaction history service.
	var start_date string = ""
	if len(os.Args) >= 4 && len(os.Args[3]) > 0 {
		start_date = os.Args[3]
	}

	var end_date string = ""
	if len(os.Args) >= 5 && len(os.Args[4]) > 0 {
		end_date = os.Args[4]
	}

	// Prepare GET request
	http_client := http.Client{
		Timeout: time.Duration(api_client.config.RequestTimeout) * time.Second,
	}
	wallet_id := os.Args[2]
	base_url := api_client.config.Server.GetURL()
	full_url := base_url + "/wallets/{" + wallet_id + "}/transaction_history" //
	started_query_string := false
	if len(start_date) > 0 {
		full_url += "?from=" + start_date
		started_query_string = true
	}
	if len(end_date) > 0 {
		if started_query_string {
			full_url += "&to=" + end_date
		} else {
			full_url += "?to=" + end_date
		}
		started_query_string = true
	}
	response, err := http_client.Get(full_url)
	if err != nil {
		fmt.Println("HTTP error occurred: ", err.Error())
		return
	}

	// Parse result
	bytes := make([]byte, response.ContentLength)
	_, err = response.Body.Read(bytes)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			fmt.Println("Error reading response: ", err.Error())
			return
		}
	}
	response_body := responses.TransactionHistory{}
	err = json.Unmarshal(bytes, &response_body)
	if err != nil {
		fmt.Println("Error parsing JSON response. ", string(bytes))
		return
	}

	// Print result to console
	fmt.Println("Request status: ", ConvertToString(response_body.Status))
	switch response_body.Status {
	case responses.Status_successful:
		if len(response_body.History) == 0 {
			fmt.Println("No transactions found")
		} else {
			fmt.Println("Transaction history")
			fmt.Println("Date (YYYYMMDD), Type (D/W), Currency, Amount")
			for _, row := range response_body.History {
				fmt.Println(row.Date, ", ", row.Type, ", ", row.Currency, ", ", row.Amount)
			}
		}
	case responses.Status_failed:
		fmt.Println("Error message: ", response_body.ErrorMessage)
	case responses.Status_unknown:
	}

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
