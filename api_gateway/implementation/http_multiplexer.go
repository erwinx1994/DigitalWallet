package implementation

import (
	"api_gateway/config"
	"api_gateway/paths"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"shared/messages"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type QueueItem struct {
	id        int64
	json_body []byte
}

type HTTPRequestMultiplexer struct {
	config    *config.Config
	context   context.Context
	global_id atomic.Int64
	// Global unique id for tracking messages. This is required if there are multiple threads
	// putting items in the queue. The responses may not arrive in order.
	// should use twitter's snowflake approach or a globally unique identifier but lack of time.
	// This is ignored for now. Assume only a single thread puts items on the queue.

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

func (mux *HTTPRequestMultiplexer) send_request_and_return_response(
	bytes_to_send []byte,
	backend_service *config.Service,
	requests_queue *redis.Client,
	responses_queue *redis.Client,
	writer http.ResponseWriter) {

	// Define aliases
	timeout := time.Duration(backend_service.RequestsQueue.Timeout) * time.Second
	queue_name := backend_service.RequestsQueue.QueueName

	// Put request in queue
	timeout_context, cancel := context.WithTimeout(mux.context, timeout)
	_, err := requests_queue.LPush(timeout_context, queue_name, bytes_to_send).Result()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		cancel()
		return
	}
	cancel()

	// Define aliases
	timeout = time.Duration(backend_service.ResponsesQueue.Timeout) * time.Second
	queue_name = backend_service.ResponsesQueue.QueueName

	// Wait for response to request
	timeout_context, cancel = context.WithTimeout(mux.context, timeout)
	string_slice, err := responses_queue.BRPop(timeout_context, timeout, queue_name).Result()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		cancel()
		return
	}
	cancel()

	// string_slice[0] gives the name of the queue
	// string_slice[1] gives the data retrieved from the queue
	// Response is prepared by the primary backend service, not the API gateway.

	// Send response to user
	writer.Write([]byte(string_slice[1]))
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

	// Prepare redis message
	body.WalletID = wallet_id
	body.Header.MessageID = mux.global_id.Add(1)
	body.Header.Action = messages.Action_deposit
	bytes, err = json.Marshal(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mux.send_request_and_return_response(
		bytes,
		&mux.config.DepositsService,
		mux.deposit_requests_queue,
		mux.deposit_responses_queue,
		writer)
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

	// Prepare redis message
	body.WalletID = wallet_id
	body.Header.MessageID = mux.global_id.Add(1)
	body.Header.Action = messages.Action_withdraw
	bytes, err = json.Marshal(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mux.send_request_and_return_response(
		bytes,
		&mux.config.WithdrawalService,
		mux.withdrawal_requests_queue,
		mux.withdrawal_responses_queue,
		writer)
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

	// Prepare redis message
	body.Header.MessageID = mux.global_id.Add(1)
	body.Header.Action = messages.Action_transfer
	bytes, err = json.Marshal(body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mux.send_request_and_return_response(
		bytes,
		&mux.config.TransferService,
		mux.transfer_requests_queue,
		mux.transfer_responses_queue,
		writer)

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

	// Prepare redis message
	request_message := messages.GET_Balance{
		Header: messages.Header{
			MessageID: mux.global_id.Add(1),
			Action:    messages.Action_get_balance,
		},
		WalletID: wallet_id,
	}
	bytes, err := json.Marshal(request_message)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mux.send_request_and_return_response(
		bytes,
		&mux.config.BalanceService,
		mux.balance_requests_queue,
		mux.balance_responses_queue,
		writer)

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

	// Prepare redis message
	request_message := messages.GET_TransactionHistory{
		Header: messages.Header{
			MessageID: mux.global_id.Add(1),
			Action:    messages.Action_get_transaction_history,
		},
		WalletID: wallet_id,
		From:     from,
		To:       to,
	}
	bytes, err := json.Marshal(request_message)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	mux.send_request_and_return_response(
		bytes,
		&mux.config.TransactionHistoryService,
		mux.transaction_history_requests_queue,
		mux.transaction_history_responses_queue,
		writer)

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
