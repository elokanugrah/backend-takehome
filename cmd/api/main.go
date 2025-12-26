package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elokanugrah/backend-takehome/internal/auth"
	"github.com/elokanugrah/backend-takehome/internal/config"
	"github.com/elokanugrah/backend-takehome/internal/database"
	"github.com/elokanugrah/backend-takehome/internal/repository"
	"github.com/elokanugrah/backend-takehome/internal/usecase"

	httpDelivery "github.com/elokanugrah/backend-takehome/internal/handler/http"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	db := database.NewMySQLDB(cfg)
	defer db.Close()

	// Initialize Auth Service
	authService := auth.NewJWTService(cfg.JWTSecret)

	// Initialize Repository Layer
	userRepo := repository.NewUserRepository(db)
	// txManager := repository.NewTransactionManager(db)

	// Initialize Usecase Layer
	userUseCase := usecase.NewUserUseCase(userRepo, authService)

	// Initialize Delivery Layer (Handler)
	apiHandler := httpDelivery.NewHandler(userUseCase, authService)

	// Setup Router
	router := httpDelivery.SetupRouter(apiHandler)

	// --- GRACEFUL SHUTDOWN SETUP ---
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Start server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Starting server on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
