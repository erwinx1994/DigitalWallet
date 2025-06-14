# shared library

The shared library contains code reused by multiple backend applications. 

To link to this library, add the following directive to the **go.mod** file in the api_client, api_gateway, balance_service, deposit_service, transaction_history_service, transfer_service and withdraw service appications.

    replace "shared" => ../shared

    