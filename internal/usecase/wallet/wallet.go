package wallet

import (
	"context"

	"github.com/google/uuid"
	
	walletErrors "github.com/totorialman/go-test-ac/internal/errors/wallet"
	"github.com/totorialman/go-test-ac/internal/repository/wallet"
	"github.com/totorialman/go-test-ac/internal/domain"
)

type Usecase struct {
	repo repository
}

func NewUsecase(repo repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Operate(ctx context.Context, w Wallet) (int64, error) {
	dbWallet := wallet.WalletDB{
		ID:     w.ID,
		Amount: w.Amount,
	}

	switch w.OperationType {
	case domain.Deposit:
		return u.repo.Deposit(ctx, dbWallet)
	case domain.Withdraw:
		return u.repo.Withdraw(ctx, dbWallet)
	default:
		return 0, walletErrors.ErrInvalidOperation
	}
}

func (u *Usecase) Balance(ctx context.Context, id uuid.UUID) (int64, error) {
	return u.repo.GetBalance(ctx, id)
}
