package wallet

import "github.com/google/uuid"

type WalletRequest struct {
	ID            uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
}

type WalletResponse struct {
	ID      uuid.UUID `json:"walletId"`
	Balance int64     `json:"balance"`
}

