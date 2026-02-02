package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/GunarsK-portfolio/portfolio-common/health"
)

type pgxChecker struct {
	pool *pgxpool.Pool
}

// NewPgxChecker returns a health.Checker that pings the pgx pool.
func NewPgxChecker(pool *pgxpool.Pool) health.Checker {
	if pool == nil {
		panic("NewPgxChecker: pool cannot be nil")
	}
	return &pgxChecker{pool: pool}
}

func (c *pgxChecker) Name() string {
	return "postgres"
}

func (c *pgxChecker) Check(ctx context.Context) health.CheckResult {
	start := time.Now()

	if err := c.pool.Ping(ctx); err != nil {
		return health.CheckResult{
			Status:  health.StatusUnhealthy,
			Latency: time.Since(start).String(),
			Error:   err.Error(),
		}
	}

	return health.CheckResult{
		Status:  health.StatusHealthy,
		Latency: time.Since(start).String(),
	}
}
