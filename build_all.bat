cd ./api_client 
go get all
go mod tidy
go build

cd ../api_gateway
go get all
go mod tidy
go build

cd ../deposit_service
go get all
go mod tidy
go build

cd ../withdraw_service
go get all
go mod tidy
go build

cd ../transfer_service
go get all
go mod tidy
go build

cd ../balance_service
go get all
go mod tidy
go build

cd ../transaction_history_service
go get all
go mod tidy
go build
