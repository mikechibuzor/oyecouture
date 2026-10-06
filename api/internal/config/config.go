package config

import (
	"net/url"
	"os"
)

const (
	defaultAPIAddress   = ":8080"
	defaultDatabaseURL  = "postgres://oye:oye@localhost:5432/oyecouture?sslmode=disable"
	defaultOwnerEmail   = "iwuagwuchibuzor@gmail.com"
	// Printed on the contact page until the owner sets one in the admin.
	defaultContactEmail = "hello@oyecouture.com"
	defaultWebBaseURL   = "http://localhost:3000"
	defaultAnnouncement = "COMPLIMENTARY TAILORING CONSULTATION · WORLDWIDE DELIVERY · BE EXCLUSIVE"
	// Origin-relative so product image URLs work in any environment without
	// knowing the host. Point R2_PUBLIC_BASE at the bucket URL to serve from R2.
	defaultR2PublicBase = "/images/catalogue"
)

type Config struct {
	APIAddress          string
	AppEnv              string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	CloudinaryCloudName string
	CloudinaryFolder    string
	DatabaseURL         string
	ContactEmail        string
	EmailFrom           string
	EmailProvider       string
	OwnerNotifyEmail    string
	PaymentProvider     string
	PaystackSecretKey   string
	R2PublicBase        string
	SendByteAPIKey      string
	WebBaseURL          string
	WhatsAppNumber      string
	Announcement        string
}

// Load resolves environment-backed application configuration with development defaults.
func Load() Config {
	return Config{
		APIAddress:          apiAddress(os.Getenv("API_BASE_URL")),
		AppEnv:              valueOrDefault(os.Getenv("APP_ENV"), "development"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryFolder:    valueOrDefault(os.Getenv("CLOUDINARY_FOLDER"), "oyecouture"),
		DatabaseURL:         valueOrDefault(os.Getenv("DATABASE_URL"), defaultDatabaseURL),
		ContactEmail:        valueOrDefault(os.Getenv("CONTACT_EMAIL"), defaultContactEmail),
		EmailFrom:           os.Getenv("EMAIL_FROM"),
		EmailProvider:       valueOrDefault(os.Getenv("EMAIL_PROVIDER"), "sendbyte"),
		OwnerNotifyEmail:    valueOrDefault(os.Getenv("OWNER_NOTIFY_EMAIL"), defaultOwnerEmail),
		PaymentProvider:     valueOrDefault(os.Getenv("PAYMENT_PROVIDER"), "paystack"),
		PaystackSecretKey:   os.Getenv("PAYSTACK_SECRET_KEY"),
		R2PublicBase:        valueOrDefault(os.Getenv("R2_PUBLIC_BASE"), defaultR2PublicBase),
		SendByteAPIKey:      os.Getenv("SENDBYTE_API_KEY"),
		WebBaseURL:          valueOrDefault(os.Getenv("WEB_BASE_URL"), defaultWebBaseURL),
		WhatsAppNumber:      os.Getenv("WHATSAPP_NUMBER"),
		Announcement:        valueOrDefault(os.Getenv("ANNOUNCEMENT_TEXT"), defaultAnnouncement),
	}
}

// apiAddress extracts the configured API port or returns the development address.
func apiAddress(apiBaseURL string) string {
	parsedURL, err := url.Parse(apiBaseURL)
	if err != nil || parsedURL.Port() == "" {
		return defaultAPIAddress
	}

	return ":" + parsedURL.Port()
}

// valueOrDefault returns a fallback only when an environment value is empty.
func valueOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}
