package responses

const (
	Status_unknown    int = 0
	Status_successful int = 1
	Status_failed     int = 2
)

type Header struct {
	MessageID int64 `json:"id"`
	Action    int   `json:"action"`
}
type Deposit struct {
	Header     Header `json:"header"`
	Status     int    `json:"status,omitempty"`
	Currency   string `json:"currency,omitempty"`
	NewBalance string `json:"new_balance,omitempty"`
}

type Withdraw struct {
	Header     Header `json:"header"`
	Status     int    `json:"status,omitempty"`
	Currency   string `json:"currency,omitempty"`
	NewBalance string `json:"new_balance,omitempty"`
}

type Transfer struct {
	Header     Header `json:"header"`
	Status     int    `json:"status,omitempty"`
	Currency   string `json:"currency,omitempty"`
	NewBalance string `json:"new_balance,omitempty"`
}

type Balance struct {
	Header   Header `json:"header"`
	Status   int    `json:"status,omitempty"`
	Currency string `json:"currency,omitempty"`
	Balance  string `json:"balance,omitempty"`
}

type Transaction struct {
	Date     string `json:"date,omitempty"`
	Type     int    `json:"type,omitempty"`
	Currency string `json:"currency,omitempty"`
	Amount   string `json:"amount,omitempty"`
}

type TransactionHistory struct {
	Header  Header        `json:"header"`
	Status  int           `json:"status,omitempty"`
	History []Transaction `json:"history,omitempty"`
}
