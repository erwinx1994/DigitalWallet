package messages

const (
	Action_unknown                 int = 0
	Action_deposit                 int = 1
	Action_withdraw                int = 2
	Action_transfer                int = 3
	Action_get_balance             int = 4
	Action_get_transaction_history int = 5
)

type RedisMessage struct {
	MessageID int64       `json:"id"`
	Action    int         `json:"action"`
	Body      interface{} `json:"body,omitempty"`
}

type POST_Deposit struct {
	WalletID string `json:"wallet_id,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type POST_Withdraw struct {
	WalletID string `json:"wallet_id,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type POST_Transfer struct {
	SourceWalletID      string `json:"source_wallet_id,omitempty"`
	DestinationWalletID string `json:"destination_wallet_id,omitempty"`
	Amount              string `json:"amount,omitempty"`
	Currency            string `json:"currency,omitempty"`
}

type GET_Balance struct {
	WalletID string `json:"wallet_id,omitempty"`
}

type GET_TransactionHistory struct {
	WalletID string `json:"wallet_id,omitempty"`
	From     string `json:"from,omitempty"` // YYYYMMDD
	To       string `json:"to,omitempty"`   // YYYYMMDD
}
