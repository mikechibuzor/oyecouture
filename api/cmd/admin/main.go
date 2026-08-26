// Command admin creates or updates the admin user for the ghost admin.
//
//	go run ./cmd/admin
//
// The password is read without echoing and stored only as a bcrypt hash. It is
// never written to .env, logs or the repository.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/term"

	"github.com/uzor-co/oyecouture/api/internal/admin"
	"github.com/uzor-co/oyecouture/api/internal/config"
)

const minimumPasswordLength = 12

func main() {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")
	applicationConfig := config.Load()

	pool, err := pgxpool.New(context.Background(), applicationConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Admin email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("read email: %v", err)
	}
	email = strings.TrimSpace(email)
	if !strings.Contains(email, "@") {
		log.Fatal("that does not look like an email address")
	}

	password, err := readPassword(reader, "Password (min 12 chars, not shown): ")
	if err != nil {
		log.Fatalf("read password: %v", err)
	}
	if len(password) < minimumPasswordLength {
		log.Fatalf("password must be at least %d characters", minimumPasswordLength)
	}

	confirmation, err := readPassword(reader, "Confirm password: ")
	if err != nil {
		log.Fatalf("read confirmation: %v", err)
	}
	if password != confirmation {
		log.Fatal("passwords do not match")
	}

	passwordHash, err := admin.HashPassword(password)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	store := admin.NewPostgresStore(pool)
	if _, err := store.CreateUser(context.Background(), email, passwordHash, "owner"); err != nil {
		log.Fatalf("save admin user: %v", err)
	}

	fmt.Printf("\nAdmin user ready: %s\nSign in at /admin/login\n", email)
}

// readPassword takes input without echoing it when attached to a terminal.
// When stdin is piped (scripted setup, CI) it falls back to reading a line,
// the password is then only as private as the calling script makes it.
func readPassword(reader *bufio.Reader, prompt string) (string, error) {
	if !term.IsTerminal(int(syscall.Stdin)) {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		return strings.TrimRight(line, "\r\n"), nil
	}

	fmt.Print(prompt)
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	return string(passwordBytes), nil
}
