package dto

type Ledger struct {
	AccountNumber string
	Amount        int64
	Idempotent    string
}
