package implementation

import (
	"api_gateway/paths"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"shared/messages"
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

	// Extract body of message
	if request.ContentLength <= 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes := make([]byte, request.ContentLength)
	_, err := request.Body.Read(bytes)
	if err != nil && !errors.Is(err, io.EOF) {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := messages.POST_Deposit{}
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify that input is correct. Basic checks only due to time limit.
	wallet_id, exist := input.WildcardSegments["wallet_id"]
	if !exist {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(wallet_id) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body.Amount) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body.Currency) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Put request in queue

	// Wait for response to request

	// Send response to user
	writer.Write(bytes)

}

func (mux *HTTPRequestMultiplexer) POST_Withdrawal(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Extract body of message
	if request.ContentLength <= 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes := make([]byte, request.ContentLength)
	_, err := request.Body.Read(bytes)
	if err != nil && !errors.Is(err, io.EOF) {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := messages.POST_Withdraw{}
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify that input is correct
	wallet_id, exist := input.WildcardSegments["wallet_id"]
	if !exist {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(wallet_id) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body.Amount) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body.Currency) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) POST_Transfer(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Extract body of message
	if request.ContentLength <= 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes := make([]byte, request.ContentLength)
	_, err := request.Body.Read(bytes)
	if err != nil && !errors.Is(err, io.EOF) {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := messages.POST_Transfer{}
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify that input is correct
	if len(body.SourceWalletID) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body.DestinationWalletID) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body.Amount) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body.Currency) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) POST_Test(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Extract body of message
	if request.ContentLength <= 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes := make([]byte, request.ContentLength)
	_, err := request.Body.Read(bytes)
	if err != nil && !errors.Is(err, io.EOF) {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := messages.POST_Deposit{}
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify that input is correct

	// Put request in queue

	// Wait for response to request

	// Send response to user
	writer.Write(bytes)
}

func (mux *HTTPRequestMultiplexer) GET_WalletBalance(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Verify that input is correct
	wallet_id, exist := input.WildcardSegments["wallet_id"]
	if !exist {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(wallet_id) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// Put request in queue

	// Wait for response to request

	// Send response to user

}

func (mux *HTTPRequestMultiplexer) GET_TransactionHistory(input *paths.MatchResult, writer http.ResponseWriter, request *http.Request) {

	// Verify that input is correct
	wallet_id, exist := input.WildcardSegments["wallet_id"]
	if !exist {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(wallet_id) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	from, exist := input.KeyValuePairs["from"]
	if !exist {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(from) != 8 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	to, exist := input.KeyValuePairs["to"]
	if exist && len(to) != 8 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

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

	// Test GET request
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

	// Test POST request
	result = paths.MatchAndExtract(request.URL.Path, paths.Test)
	if result.MatchFound {
		mux.POST_Test(result, writer, request)
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
