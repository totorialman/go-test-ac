package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustInitDB(ctx context.Context) *pgxpool.Pool {
	dsn, err := LoadConfigDB()
	if err != nil {
		log.Fatalf("failed to load DB config: %v", err)
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("failed to parse DB config: %v", err)
	}

	cfg.MaxConns = 100
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create DB pool: %v", err)
	}

	for i := range 5 {
		select {
		case <-ctx.Done():
			log.Fatalf("DB init cancelled: %v", ctx.Err())
		default:
		}

		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		err := pool.Ping(pingCtx)
		cancel()

		if err == nil {
			log.Println("DB connection established")
			return pool
		}

		log.Printf("DB not ready (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("DB connection failed after retries")
	return nil
}
