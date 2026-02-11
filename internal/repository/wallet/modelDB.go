package wallet

import "github.com/google/uuid"

type WalletDB struct {
	ID     uuid.UUID
	Amount int64
}
