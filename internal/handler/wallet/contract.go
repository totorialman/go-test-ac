package wallet

import (
	"context"

	"github.com/google/uuid"

	"github.com/totorialman/go-test-ac/internal/usecase/wallet"
)

type usecase interface {
	Operate(ctx context.Context, w wallet.Wallet) (int64, error)
	Balance(ctx context.Context, id uuid.UUID) (int64, error)
}
