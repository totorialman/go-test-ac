//go:generate mockgen -source=contract.go -destination=wallet_repository_mocks_test.go -package=wallet_test
package wallet

import (
	"context"

	"github.com/google/uuid"

	"github.com/totorialman/go-test-ac/internal/repository/wallet"
)

type repository interface {
	GetBalance(ctx context.Context, id uuid.UUID) (int64, error)
	Deposit(ctx context.Context, w wallet.WalletDB) (int64, error)
	Withdraw(ctx context.Context, w wallet.WalletDB) (int64, error)
}
