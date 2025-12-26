package main

import (
	"log"

	"github.com/elokanugrah/backend-takehome/app/config"
	"github.com/elokanugrah/backend-takehome/app/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	db, err := database.NewMySQLDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Starting database seeding...")

	// Define users to seed
	users := []struct {
		Name     string
		Email    string
		Password string
	}{
		{"Alice Wonderland", "alice@example.com", "password"},
		{"Bob Builder", "bob@example.com", "password"},
	}

	query := "INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)"

	for _, u := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", u.Email, err)
			continue
		}

		if _, err := db.Exec(query, u.Name, u.Email, string(hashedPassword)); err != nil {
			log.Printf("Failed to insert user %s: %v", u.Email, err)
		} else {
			log.Printf("Inserted user: %s", u.Email)
		}
	}

	log.Println("Seeding completed.")
}
