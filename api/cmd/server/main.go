package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/uzor-co/oyecouture/api/internal/config"
	"github.com/uzor-co/oyecouture/api/internal/server"
	"github.com/uzor-co/oyecouture/api/migrations"
)

// main prepares infrastructure, starts the API, and coordinates graceful shutdown.
func main() {
	// Load local development values before resolving the application configuration.
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	applicationConfig := config.Load()

	// Bring the database schema and seed data to the version required by this binary.
	if err := migrations.Up(applicationConfig.DatabaseURL); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	// Create the shared application pool after migrations release their connection.
	databasePool, err := pgxpool.New(context.Background(), applicationConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer databasePool.Close()

	httpServer := server.New(applicationConfig, databasePool)
	serverErrors := make(chan error, 1)

	// Serve requests asynchronously while the main goroutine monitors exit conditions.
	go func() {
		log.Printf("api listening on %s", httpServer.Addr)
		serverErrors <- httpServer.ListenAndServe()
	}()

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

	// Wait until the process receives a shutdown signal or the HTTP server stops.
	select {
	case shutdownSignal := <-shutdownSignals:
		log.Printf("received %s, shutting down", shutdownSignal)
	case serverError := <-serverErrors:
		if !errors.Is(serverError, http.ErrServerClosed) {
			log.Fatalf("serve api: %v", serverError)
		}
		return
	}

	// Give in-flight requests a bounded window to complete before process exit.
	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		log.Printf("shutdown api: %v", err)
	}
}
