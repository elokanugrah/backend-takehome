package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elokanugrah/backend-takehome/internal/usecase"
)

// txKey is the key used to store the transaction object in the context.
type txKey struct{}

type SqlTransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) usecase.TransactionManager {
	return &SqlTransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction.
// It begins a transaction, calls the provided function with a new context
// containing the transaction, and then commits or rolls back based on the error.
func (tm *SqlTransactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a new context with the transaction object
	txCtx := context.WithValue(ctx, txKey{}, tx)

	// Use a deferred function to recover from panics and rollback the transaction.
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-panic after rolling back
		}
	}()

	// Execute the provided function with the transaction context
	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, unable to rollback: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// dbExecutor defines the methods shared by *sql.DB and *sql.Tx
type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
