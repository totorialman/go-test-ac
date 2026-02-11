package wallet

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	walletErrors "github.com/totorialman/go-test-ac/internal/errors/wallet"
    "github.com/totorialman/go-test-ac/internal/usecase/wallet"
    "github.com/totorialman/go-test-ac/internal/domain"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) Operate(w http.ResponseWriter, r *http.Request) {
	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("decode error: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		log.Printf("invalid amount: %d", req.Amount)
		http.Error(w, walletErrors.ErrInvalidAmount.Error(), http.StatusBadRequest)
		return
	}
	if req.OperationType != domain.Deposit && req.OperationType != domain.Withdraw {
		log.Printf("invalid operation type: %s", req.OperationType)
		http.Error(w, walletErrors.ErrInvalidOperation.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("operate request: id=%s type=%s amount=%d",
		req.ID, req.OperationType, req.Amount,
	)

	wallet := wallet.Wallet{
		ID:            req.ID,
		OperationType: req.OperationType,
		Amount:        req.Amount,
	}

	newBalance, err := h.usecase.Operate(r.Context(), wallet)
	if err != nil {
		log.Printf("operate error: %v", err)

		switch {
		case errors.Is(err, walletErrors.ErrNotEnoughFunds):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, walletErrors.ErrWalletNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, walletErrors.ErrInvalidAmount),
			errors.Is(err, walletErrors.ErrInvalidOperation):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("operate success: id=%s new_balance=%d", wallet.ID, newBalance)

	res := WalletResponse{
		ID:      wallet.ID,
		Balance: newBalance,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("JSON encode error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["WALLET_UUID"]
	log.Printf("balance request: id=%s", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("invalid uuid: %v", err)
		http.Error(w, "invalid wallet id", http.StatusBadRequest)
		return
	}

	balance, err := h.usecase.Balance(r.Context(), id)
	if err != nil {
		log.Printf("balance error: %v", err)

		if errors.Is(err, walletErrors.ErrWalletNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("balance success: id=%s balance=%d", id, balance)

	res := WalletResponse{
		ID:      id,
		Balance: balance,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("JSON encode error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
