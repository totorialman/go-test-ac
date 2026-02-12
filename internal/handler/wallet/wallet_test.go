package wallet_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/totorialman/go-test-ac/internal/domain"
	walletErrors "github.com/totorialman/go-test-ac/internal/errors/wallet"
	"github.com/totorialman/go-test-ac/internal/handler/wallet"
)

func TestHandler_Operate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockusecase(ctrl)
	h := wallet.NewHandler(mockUsecase)

	walletID := uuid.New()

	tests := []struct {
		name           string
		reqBody        wallet.WalletRequest
		mockReturn     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid deposit",
			reqBody: wallet.WalletRequest{
				ID:            walletID,
				OperationType: domain.Deposit,
				Amount:        500,
			},
			mockReturn: func() {
				mockUsecase.EXPECT().
					Operate(gomock.Any(), gomock.Any()).
					Return(int64(1500), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"balance":1500`,
		},
		{
			name: "invalid amount",
			reqBody: wallet.WalletRequest{
				ID:            walletID,
				OperationType: domain.Deposit,
				Amount:        -10,
			},
			mockReturn:     func() {}, 
			expectedStatus: http.StatusBadRequest,
			expectedBody:   walletErrors.ErrInvalidAmount.Error(),
		},
		{
			name: "not enough funds",
			reqBody: wallet.WalletRequest{
				ID:            walletID,
				OperationType: domain.Withdraw,
				Amount:        1000,
			},
			mockReturn: func() {
				mockUsecase.EXPECT().
					Operate(gomock.Any(), gomock.Any()).
					Return(int64(0), walletErrors.ErrNotEnoughFunds)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   walletErrors.ErrNotEnoughFunds.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReturn()

			bodyBytes, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/operate", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			h.Operate(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			assert.Contains(t, buf.String(), tt.expectedBody)
		})
	}
}

func TestHandler_Balance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockusecase(ctrl)
	h := wallet.NewHandler(mockUsecase)

	validID := uuid.New()
	invalidID := "not-a-uuid"

	tests := []struct {
		name           string
		walletID       string
		mockReturn     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "valid wallet",
			walletID: validID.String(),
			mockReturn: func() {
				mockUsecase.EXPECT().
					Balance(gomock.Any(), validID).
					Return(int64(2000), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"balance":2000`,
		},
		{
			name:           "invalid uuid",
			walletID:       invalidID,
			mockReturn:     func() {}, 
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid wallet id",
		},
		{
			name:     "wallet not found",
			walletID: validID.String(),
			mockReturn: func() {
				mockUsecase.EXPECT().
					Balance(gomock.Any(), validID).
					Return(int64(0), walletErrors.ErrWalletNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   walletErrors.ErrWalletNotFound.Error(),
		},
		{
			name:     "internal server error",
			walletID: validID.String(),
			mockReturn: func() {
				mockUsecase.EXPECT().
					Balance(gomock.Any(), validID).
					Return(int64(0), errors.New("some internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReturn()

			req := httptest.NewRequest(http.MethodGet, "/balance", nil)
			req = mux.SetURLVars(req, map[string]string{
				"WALLET_UUID": tt.walletID,
			})
			w := httptest.NewRecorder()

			h.Balance(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			assert.Contains(t, buf.String(), tt.expectedBody)
		})
	}
}
