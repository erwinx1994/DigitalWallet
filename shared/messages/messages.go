package messages

const (
	Action_unknown                 int = 0
	Action_deposit                 int = 1
	Action_withdraw                int = 2
	Action_transfer                int = 3
	Action_get_balance             int = 4
	Action_get_transaction_history int = 5
)

type Header struct {
	MessageID int64 `json:"id"`
	Action    int   `json:"action"`
}

type POST_Deposit struct {
	Header   Header `json:"header"`
	WalletID string `json:"wallet_id,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type POST_Withdraw struct {
	Header   Header `json:"header"`
	WalletID string `json:"wallet_id,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type POST_Transfer struct {
	Header              Header `json:"header"`
	SourceWalletID      string `json:"source_wallet_id,omitempty"`
	DestinationWalletID string `json:"destination_wallet_id,omitempty"`
	Amount              string `json:"amount,omitempty"`
	Currency            string `json:"currency,omitempty"`
}

type GET_Balance struct {
	Header   Header `json:"header"`
	WalletID string `json:"wallet_id,omitempty"`
}

type GET_TransactionHistory struct {
	Header   Header `json:"header"`
	WalletID string `json:"wallet_id,omitempty"`
	From     string `json:"from,omitempty"` // YYYYMMDD
	To       string `json:"to,omitempty"`   // YYYYMMDD
}
