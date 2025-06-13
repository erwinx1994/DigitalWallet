
# Digital wallet

This repository provides multiple applications that work together to function as a backend for a digital wallet service. The purpose of each application is described below.

    api_gateway
    Route incoming HTTP requests from the user to the correct backend service. 
    Return the result after processing to the user.

    deposit_service
    Deposit money into wallet

    withdraw_service
    Withdraw money from a wallet

    transfer_service
    Transfer money from one wallet to another

    balance_service
    Retrieve wallet balance

    transaction_history_service
    Retrieve transaction history of wallet

The system design diagram that connects them all can be found in the [./system_design/Simplified digital wallet system.pdf](./system_design/Simplified%20digital%20wallet%20system.pdf) file in this repository.

## Design of RESTful API

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
    GET /wallets/{wallet_id}/transaction_history

## Set up of databases 

PostgreSQL must be installed and set up correctly on your test computer before running this program.

## Set up of message queues

Redis must be installed and set up correctly on your test computer before running this program.

## How to compile the code

The go compiler is required to compile the programs in this repository it can be downloaded from [https://go.dev/dl/](https://go.dev/dl/). 

The following commands can be used to update the dependent libraries and compile each program.

    api_gateway
    cd ./api_gateway
    go get all
    go mod tidy
    go build

    deposit_service
    cd ./deposit_service
    go get all
    go mod tidy
    go build

    withdraw_service
    cd ./withdraw_service
    go get all
    go mod tidy
    go build

    transfer_service
    cd ./transfer_service
    go get all
    go mod tidy
    go build

    balance_service
    cd ./balance_service
    go get all
    go mod tidy
    go build

    transaction_history_service
    cd ./transaction_history_service
    go get all
    go mod tidy
    go build

## Running the program


## Running unit tests




