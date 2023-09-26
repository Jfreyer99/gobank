package types

type CreateTransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}
