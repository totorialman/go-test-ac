package wallet_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/totorialman/go-test-ac/internal/domain"
	wErr "github.com/totorialman/go-test-ac/internal/errors/wallet"
	repo "github.com/totorialman/go-test-ac/internal/repository/wallet"
	w "github.com/totorialman/go-test-ac/internal/usecase/wallet"
)

func TestUsecase_Operate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockrepository(ctrl)
	usecase := w.NewUsecase(mockRepo)

	userID := uuid.New()

	tests := []struct {
		name        string
		wallet      w.Wallet
		mockSetup   func()
		wantBalance int64
		wantErr     error
	}{
		{
			name: "Deposit success",
			wallet: w.Wallet{
				ID:            userID,
				OperationType: domain.Deposit,
				Amount:        100,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					Deposit(gomock.Any(), repo.WalletDB{ID: userID, Amount: 100}).
					Return(int64(150), nil)
			},
			wantBalance: 150,
			wantErr:     nil,
		},
		{
			name: "Withdraw success",
			wallet: w.Wallet{
				ID:            userID,
				OperationType: domain.Withdraw,
				Amount:        50,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					Withdraw(gomock.Any(), repo.WalletDB{ID: userID, Amount: 50}).
					Return(int64(50), nil)
			},
			wantBalance: 50,
			wantErr:     nil,
		},
		{
			name: "Invalid operation",
			wallet: w.Wallet{
				ID:            userID,
				OperationType: "invalid",
				Amount:        10,
			},
			mockSetup:   func() {},
			wantBalance: 0,
			wantErr:     wErr.ErrInvalidOperation,
		},
		{
			name: "Withdraw not enough funds",
			wallet: w.Wallet{
				ID:            userID,
				OperationType: domain.Withdraw,
				Amount:        200,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					Withdraw(gomock.Any(), repo.WalletDB{ID: userID, Amount: 200}).
					Return(int64(0), wErr.ErrNotEnoughFunds)
			},
			wantBalance: 0,
			wantErr:     wErr.ErrNotEnoughFunds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			balance, err := usecase.Operate(context.Background(), tt.wallet)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBalance, balance)
			}
		})
	}
}

func TestUsecase_Balance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockrepository(ctrl)
	usecase := w.NewUsecase(mockRepo)

	userID := uuid.New()

	tests := []struct {
		name        string
		mockSetup   func()
		wantBalance int64
		wantErr     error
	}{
		{
			name: "Get balance success",
			mockSetup: func() {
				mockRepo.EXPECT().GetBalance(gomock.Any(), userID).Return(int64(100), nil)
			},
			wantBalance: 100,
			wantErr:     nil,
		},
		{
			name: "Wallet not found",
			mockSetup: func() {
				mockRepo.EXPECT().GetBalance(gomock.Any(), userID).Return(int64(0), wErr.ErrWalletNotFound)
			},
			wantBalance: 0,
			wantErr:     wErr.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			balance, err := usecase.Balance(context.Background(), userID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBalance, balance)
			}
		})
	}
}
