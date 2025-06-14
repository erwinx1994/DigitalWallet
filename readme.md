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

## Client application

The client application can be found in the **api_client** subfolder of this repository. Once compiled, it can be used to interact with the backend applications to manage your wallet.

This command allows you to deposit an amount of money in the specified wallet in the specified currency.

	api_client deposit <wallet_id> <currency> <amount>

This command allows you to withdraw an amount of money from the specified wallet in the specified currency.

	api_client withdraw <wallet_id> <currency> <amount>

This command allows you to transfer an amount of money from source to destination wallets of the same currency.

	api_client transfer <source_wallet_id> <destination_wallet_id> <currency> <amount>

This command allows you to get the balance of the specified wallet.

	api_client get_balance <wallet_id>

This command allows you to get the transaction history of the specified wallet from a start date to end date, inclusive.

	api_client get_transaction_history <wallet_id> <start_date> <end_date>

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
    GET /wallets/{wallet_id}/transaction_history?from=YYYYMMDD&to=YYYYMMDD

## Set up of environment

Firewalls in your test environment may need to be disabled for the Digital Wallet backend to work correctly. The screenshots below show how to disable Windows firewall if you are testing on a Windows machine. 

![](./images/firewall_and_network_protection_1.png)

![](./images/firewall_and_network_protection_2.png)

![](./images/firewall_and_network_protection_3.png)

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




