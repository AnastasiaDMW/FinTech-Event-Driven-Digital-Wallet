package model

type ChangedUserEvent struct {
	UserID int    `json:"id"`
	Email  string `json:"email"`
}

type TransactionEvent struct {
	UserID int     `json:"id"`
	Status string  `json:"status"`
	Amount float64 `json:"amount"`
}
