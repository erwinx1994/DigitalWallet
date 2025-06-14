package responses

const (
	Status_unknown    int = 0
	Status_successful int = 1
)

type Deposit struct {
	Status     int    `json:"status,omitempty"`
	NewBalance string `json:"new_balance,omitempty"`
}

type Withdraw struct {
	Status     int    `json:"status,omitempty"`
	NewBalance string `json:"new_balance,omitempty"`
}

type Transfer struct {
	Status     int    `json:"status,omitempty"`
	NewBalance string `json:"new_balance,omitempty"`
}

type Balance struct {
	Status  int    `json:"status,omitempty"`
	Balance string `json:"balance,omitempty"`
}

type Transaction struct {
	Date     string `json:"date,omitempty"`
	Action   int    `json:"action,omitempty"`
	Currency string `json:"currency,omitempty"`
	Amount   string `json:"amount,omitempty"`
}

type TransactionHistory struct {
	Status  int           `json:"status,omitempty"`
	History []Transaction `json:"history,omitempty"`
}
