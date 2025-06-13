package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	content_type_json string = "application/json"
	http_timeout      int    = 30 // s
)

/*
The following HTTP requests need to be tested.

Deposit money into wallet
POST /wallets/{wallet_id}/deposits

	{
	    "amount": 5000,
	    "currency": "XXX"
	}

Withdraw money from a wallet
POST /wallets/{wallet_id}/withdrawals

	{
	    "amount": 5000,
	    "currency": "XXX"
	}

Transfer money from one wallet to another
POST /transfer

	{
	    "source_wallet_id": "id1",
	    "destination_wallet_id": "id2",
	    "amount": 5000,
	    "currency": "XXX"
	}

Retrieve balance of a wallet
GET /wallets/{wallet_id}/balance

Retrieve transaction history of a wallet
GET /wallets/{wallet_id}/transaction_history?from=YYYYMMDD&to=YYYYMMDD
*/
func test_http_client() {

	http_client := http.Client{
		Timeout: time.Duration(http_timeout) * time.Second,
	}

	get_balance := http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:1120",
			Path:   "/test",
		},
	}

	response, err := http_client.Do(&get_balance)
	if err != nil {
		log.Fatal("HTTP GET request failed.", err.Error())
	}

	fmt.Println(response.Header)
	bytes_read := make([]byte, response.ContentLength)
	number_of_bytes_read, err := response.Body.Read(bytes_read)
	fmt.Println(number_of_bytes_read)
	for number_of_bytes_read > 0 {
		fmt.Println(string(bytes_read[:number_of_bytes_read]))
		number_of_bytes_read, err = response.Body.Read(bytes_read)
	}
	defer response.Body.Close()

}

func main() {
	test_http_client()
}
