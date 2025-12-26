package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/elokanugrah/backend-takehome/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLDB initializes a new MySQL database connection
func NewMySQLDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("FATAL: Could not prepare database connection: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection is actually alive
	if err := db.Ping(); err != nil {
		log.Fatalf("FATAL: Database is not reachable: %v", err)
	}

	log.Println("Database connection successful.")
	return db
}
