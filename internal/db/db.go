package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gk-dev10/sheguard_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool
var Queries *sqlc.Queries

func Init(ctx context.Context) error {
	url := os.Getenv("SUPABASE_DB_URL")
	if url == "" {
		return errors.New("SUPABASE_DB_URL is not set")
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return err
	}

	// pool tuning - Prod level
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctxPing); err != nil {
		return err
	}

	Pool = pool
	Queries = sqlc.New(pool)
	fmt.Println("\nHurraayyy, Supabase Database Connected Successfully")
	return nil
}

// to shut down the pool.
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
