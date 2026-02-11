package wallet

import "errors"

var (
	ErrNotEnoughFunds   = errors.New("not enough funds")
	ErrWalletNotFound   = errors.New("wallet not found")
	ErrInvalidOperation = errors.New("invalid operation type")
	ErrInvalidAmount    = errors.New("amount must be positive")
)
