module transaction_history_service

go 1.24.4

replace shared => ../shared

require (
	github.com/redis/go-redis/v9 v9.10.0
	gopkg.in/yaml.v3 v3.0.1
	shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/lib/pq v1.10.9 // indirect
)
