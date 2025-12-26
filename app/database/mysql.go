package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/elokanugrah/backend-takehome/app/config"
	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLDB initializes a new MySQL database connection
func NewMySQLDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection is actually alive
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
