package implementation

import (
	"api_gateway/paths"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/redis/go-redis/v9"
)

type QueueItem struct {
	id        int64
	json_body []byte
}

type HTTPRequestMultiplexer struct {
	api_responses sync.Map

	// Resource handles to request and response queues
	deposit_requests_queue              *redis.Client
	deposit_responses_queue             *redis.Client
	withdrawal_requests_queue           *redis.Client
	withdrawal_responses_queue          *redis.Client
	transfer_requests_queue             *redis.Client
	transfer_responses_queue            *redis.Client
	balance_requests_queue              *redis.Client
	balance_responses_queue             *redis.Client
	transaction_history_requests_queue  *redis.Client
	transaction_history_responses_queue *redis.Client
}

func (mux *HTTPRequestMultiplexer) POST_Deposit(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Verify that the input is correct

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) POST_Withdrawal(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) POST_Transfer(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) GET_WalletBalance(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) GET_TransactionHistory(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) GET_Test(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Log content of request
	log.Println("Method: ", request.Method)
	log.Println("Scheme: ", request.URL.Scheme)
	log.Println("Host: ", request.URL.Host)
	log.Println("Path: ", request.URL.Path)
	log.Println("Header: ", request.Header)

	// Send test response to user
	type TestResponse struct {
		Message string `json:",omitempty"`
	}
	test_response := TestResponse{
		Message: "Hi. This is a test response.",
	}
	bytes, err := json.Marshal(test_response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write([]byte(bytes))
}

func (mux *HTTPRequestMultiplexer) ProcessGETRequests(writer http.ResponseWriter, request *http.Request) {

	// Request for balance
	result := paths.MatchAndExtract(request.URL.Path, paths.Wallets_balance)
	if result.MatchFound {
		mux.GET_WalletBalance(result, writer, request)
		return
	}

	// Request for transaction history
	result = paths.MatchAndExtract(request.URL.Path, paths.Wallets_transaction_history)
	if result.MatchFound {
		mux.GET_TransactionHistory(result, writer, request)
		return
	}

	// Test request
	result = paths.MatchAndExtract(request.URL.Path, paths.Test)
	if result.MatchFound {
		mux.GET_Test(result, writer, request)
		return
	}
}

func (mux *HTTPRequestMultiplexer) ProcessPOSTRequests(writer http.ResponseWriter, request *http.Request) {

	// Deposit money into wallet
	result := paths.MatchAndExtract(request.URL.Path, paths.Wallets_deposits)
	if result.MatchFound {
		mux.POST_Deposit(result, writer, request)
		return
	}

	// Withdraw money from wallet
	result = paths.MatchAndExtract(request.URL.Path, paths.Wallets_withdrawals)
	if result.MatchFound {
		mux.POST_Withdrawal(result, writer, request)
		return
	}

	// Transfer money from one wallet to another
	result = paths.MatchAndExtract(request.URL.Path, paths.Transfer)
	if result.MatchFound {
		mux.POST_Transfer(result, writer, request)
		return
	}

}

// A new goroutine is created to serve each HTTP request
func (mux *HTTPRequestMultiplexer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		mux.ProcessGETRequests(writer, request)
	} else if request.Method == http.MethodPost {
		mux.ProcessPOSTRequests(writer, request)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}
