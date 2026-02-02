package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/GunarsK-portfolio/portfolio-common/config"
)

// NewPool creates a pgx connection pool from the shared DatabaseConfig.
func NewPool(ctx context.Context, dbCfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	sslMode := dbCfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	connURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name, sslMode,
	)

	poolCfg, err := pgxpool.ParseConfig(connURL)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	poolCfg.MaxConns = 25
	poolCfg.MinConns = 5
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 10 * time.Minute
	poolCfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
