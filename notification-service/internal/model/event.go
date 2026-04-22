package model

type ChangedUserEvent struct {
	UserID int    `json:"id"`
	Email  string `json:"email"`
}

type TransactionEvent struct {
	ID     int     `json:"id"`
	UserID int     `json:"userId"`
	Status string  `json:"status"`
	Amount float64 `json:"amount"`
}
