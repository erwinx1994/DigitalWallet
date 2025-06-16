package responses

const (
	Status_unknown    int = 0
	Status_successful int = 1
	Status_failed     int = 2
)

const (
	Transaction_type_unknown  string = "U"
	Transaction_type_deposit  string = "D"
	Transaction_type_withdraw string = "W"
)

type Header struct {
	MessageID int64 `json:"id"`
	Action    int   `json:"action"`
}

type Deposit struct {
	Header       Header `json:"header"`
	Status       int    `json:"status,omitempty"`
	ErrorMessage string `json:"error_message,omit_empty"`
	Currency     string `json:"currency,omitempty"`
	NewBalance   string `json:"new_balance,omitempty"`
}

type Withdraw struct {
	Header       Header `json:"header"`
	Status       int    `json:"status,omitempty"`
	ErrorMessage string `json:"error_message,omit_empty"`
	Currency     string `json:"currency,omitempty"`
	NewBalance   string `json:"new_balance,omitempty"`
}

type Transfer struct {
	Header       Header `json:"header"`
	Status       int    `json:"status,omitempty"`
	ErrorMessage string `json:"error_message,omit_empty"`
	Currency     string `json:"currency,omitempty"`
	NewBalance   string `json:"new_balance,omitempty"`
}

type Balance struct {
	Header       Header `json:"header"`
	Status       int    `json:"status,omitempty"`
	ErrorMessage string `json:"error_message,omit_empty"`
	Currency     string `json:"currency,omitempty"`
	Balance      string `json:"balance,omitempty"`
}

type Transaction struct {
	Date     string `json:"date,omitempty"`
	Type     string `json:"type,omitempty"`
	Currency string `json:"currency,omitempty"`
	Amount   string `json:"amount,omitempty"`
}

type TransactionHistory struct {
	Header       Header        `json:"header"`
	Status       int           `json:"status,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
	History      []Transaction `json:"history,omitempty"`
}
