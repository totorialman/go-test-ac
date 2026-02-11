package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/totorialman/go-test-ac/internal/config"
	walletHandler "github.com/totorialman/go-test-ac/internal/handler/wallet"
	walletRepository "github.com/totorialman/go-test-ac/internal/repository/wallet"
	walletUsecase "github.com/totorialman/go-test-ac/internal/usecase/wallet"
)

func main() {
	servPort := ":" + os.Getenv("PORT")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPool := config.MustInitDB(ctx)
	defer dbPool.Close()

	walletRepo := walletRepository.NewRepository(dbPool)
	walletUsecase := walletUsecase.NewUsecase(walletRepo)
	walletHandler := walletHandler.NewHandler(walletUsecase)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/api/v1/wallet", walletHandler.Operate).Methods("POST")
	r.HandleFunc("/api/v1/wallets/{WALLET_UUID}", walletHandler.Balance).Methods("GET")

	server := &http.Server{
		Addr:         servPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server started on %s\n", servPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAnd Serve error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped gracefully")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("started %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)

		next.ServeHTTP(w, r)

		log.Printf("completed %s %s in %v\n\n", r.Method, r.RequestURI, time.Since(start))
	})
}
