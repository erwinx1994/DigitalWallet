cd ./api_client
go clean --testcache
go test ./...

cd ../api_gateway
go clean --testcache
go test ./...

cd ../balance_service
go clean --testcache
go test ./...

cd ../deposit_service
go clean --testcache
go test ./...

cd ../shared
go clean --testcache
go test ./...

cd ../transaction_history_service
go clean --testcache
go test ./...

cd ../transfer_service
go clean --testcache
go test ./...

cd ../withdraw_service
go clean --testcache
go test ./...

cd ..
