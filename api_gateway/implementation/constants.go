package implementation

// Constants are used for quick comparisons when processing responses
const (
	service_balance             int = 1
	service_deposit             int = 2
	service_transaction_history int = 3
	service_transfer            int = 4
	service_withdraw            int = 5
)

// Names are used for printing to logs
func get_service_name(service_type int) string {
	var service_name string = ""
	switch service_type {
	case service_balance:
		service_name = "balance service"
	case service_deposit:
		service_name = "deposit service"
	case service_transaction_history:
		service_name = "transaction history service"
	case service_transfer:
		service_name = "transfer service"
	case service_withdraw:
		service_name = "withdraw service"
	}
	return service_name
}
