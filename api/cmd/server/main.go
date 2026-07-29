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
)

func main() {
	_ = godotenv.Load(".env", "../.env")

	applicationConfig := config.Load()
	databasePool, err := pgxpool.New(context.Background(), applicationConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer databasePool.Close()

	httpServer := server.New(applicationConfig, databasePool)
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("api listening on %s", httpServer.Addr)
		serverErrors <- httpServer.ListenAndServe()
	}()

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case shutdownSignal := <-shutdownSignals:
		log.Printf("received %s, shutting down", shutdownSignal)
	case serverError := <-serverErrors:
		if !errors.Is(serverError, http.ErrServerClosed) {
			log.Fatalf("serve api: %v", serverError)
		}
		return
	}

	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		log.Printf("shutdown api: %v", err)
	}
}
