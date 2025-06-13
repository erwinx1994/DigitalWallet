
# Digital wallet

This repository provides multiple applications that work together to function as a backend for a digital wallet service. The purpose of each application is described below.

    api_gateway
    Route incoming HTTP requests from the user to the correct backend service. Return the result after processing to the user.

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

# How to compile the code

The following software is required to test and run the programs in this repository.

    git
    PostGre SQL
    Redis
    Go 

# Set up of databases 

# Set up of message queues


# Running the program


# Running unit tests




