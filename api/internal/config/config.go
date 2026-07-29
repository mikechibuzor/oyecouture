package config

import (
	"net/url"
	"os"
)

const (
	defaultAPIAddress  = ":8080"
	defaultDatabaseURL = "postgres://oye:oye@localhost:5432/oyecouture?sslmode=disable"
	defaultWebBaseURL  = "http://localhost:3000"
)

type Config struct {
	APIAddress  string
	AppEnv      string
	DatabaseURL string
	WebBaseURL  string
}

func Load() Config {
	return Config{
		APIAddress:  apiAddress(os.Getenv("API_BASE_URL")),
		AppEnv:      valueOrDefault(os.Getenv("APP_ENV"), "development"),
		DatabaseURL: valueOrDefault(os.Getenv("DATABASE_URL"), defaultDatabaseURL),
		WebBaseURL:  valueOrDefault(os.Getenv("WEB_BASE_URL"), defaultWebBaseURL),
	}
}

func apiAddress(apiBaseURL string) string {
	parsedURL, err := url.Parse(apiBaseURL)
	if err != nil || parsedURL.Port() == "" {
		return defaultAPIAddress
	}

	return ":" + parsedURL.Port()
}

func valueOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}
