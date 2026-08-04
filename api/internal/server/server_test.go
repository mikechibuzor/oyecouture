package server

import "testing"

func TestIsAllowedOriginPermitsLoopbackInDevelopment(t *testing.T) {
	if !isAllowedOrigin("http://127.0.0.1:3010", "http://localhost:3000", "development") {
		t.Fatal("expected loopback origin to be allowed in development")
	}
	if !isAllowedOrigin("http://localhost:5173", "http://localhost:3000", "development") {
		t.Fatal("expected localhost origin to be allowed in development")
	}
}

func TestIsAllowedOriginRequiresConfiguredOriginInProduction(t *testing.T) {
	if isAllowedOrigin("http://127.0.0.1:3010", "https://oyecouture.com", "production") {
		t.Fatal("did not expect loopback origin to be allowed in production")
	}
	if !isAllowedOrigin("https://oyecouture.com", "https://oyecouture.com", "production") {
		t.Fatal("expected configured production origin to be allowed")
	}
}
