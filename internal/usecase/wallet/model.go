package wallet

import "github.com/google/uuid"

type Wallet struct {
	ID            uuid.UUID
	OperationType string
	Amount        int64
}
