package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthContext carries session identity for audit.set_context() calls.
type AuthContext struct {
	UserID    int64
	Username  string
	ClientIP  string
	UserAgent string
}

// Repository defines the interface for data access.
type Repository interface {
	Ping(ctx context.Context) error
}

type repository struct {
	pool *pgxpool.Pool
}

// New creates a new repository backed by a pgx pool.
func New(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

func (r *repository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

// withAuditTx begins a transaction, sets audit context, executes fn, and commits.
//
//nolint:unused
func (r *repository) withAuditTx(ctx context.Context, auth AuthContext, fn func(tx pgx.Tx) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx,
		"SELECT audit.set_context($1, $2, $3, $4)",
		auth.UserID, auth.Username, auth.ClientIP, auth.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("set audit context: %w", err)
	}

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// callFunctionInTx executes a query inside a transaction with audit context set.
// Returns the JSONB result from the database function.
//
//nolint:unused
func (r *repository) callFunctionInTx(ctx context.Context, auth AuthContext, query string, args ...any) (json.RawMessage, error) {
	var result json.RawMessage
	err := r.withAuditTx(ctx, auth, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, args...).Scan(&result)
	})
	return result, err
}

// execFunctionInTx executes a query inside a transaction with audit context set.
// Returns the boolean result from delete functions.
//
//nolint:unused
func (r *repository) execFunctionInTx(ctx context.Context, auth AuthContext, query string, args ...any) (bool, error) {
	var result bool
	err := r.withAuditTx(ctx, auth, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, query, args...).Scan(&result)
	})
	return result, err
}
