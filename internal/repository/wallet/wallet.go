package wallet

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/totorialman/go-test-ac/internal/errors/wallet"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	var balance int64
	err := r.db.QueryRow(ctx, `SELECT balance FROM wallets WHERE id = $1`, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, wallet.ErrWalletNotFound
		}
		return 0, err
	}
	return balance, nil
}

func (r *Repository) Deposit(ctx context.Context, w WalletDB) (int64, error) {
	var newBalance int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO wallets (id, balance)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET balance = wallets.balance + EXCLUDED.balance
		RETURNING balance
	`, w.ID, w.Amount).Scan(&newBalance)
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}

func (r *Repository) Withdraw(ctx context.Context, w WalletDB) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var currentBalance int64
	err = tx.QueryRow(ctx, `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`, w.ID).Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, wallet.ErrWalletNotFound
		}
		return 0, err
	}

	if currentBalance < w.Amount {
		return 0, wallet.ErrNotEnoughFunds
	}

	var newBalance int64
	err = tx.QueryRow(ctx, `UPDATE wallets SET balance = $2 WHERE id = $1 RETURNING balance`, w.ID, currentBalance-w.Amount).Scan(&newBalance)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return newBalance, nil
}
