package server

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/uzor-co/oyecouture/api/internal/config"
	"github.com/uzor-co/oyecouture/api/internal/health"
	"github.com/uzor-co/oyecouture/api/internal/orders"
	"github.com/uzor-co/oyecouture/api/internal/products"
)

// New wires infrastructure, services, and handlers into the API server.
func New(applicationConfig config.Config, databasePool *pgxpool.Pool) *http.Server {
	// Build the product read stack once and expose only its public storefront routes.
	productRepository := products.NewPostgresRepository(databasePool)
	productService := products.NewService(productRepository, applicationConfig.R2PublicBase)
	productHandler := products.NewHandler(productService)

	// Orders are priced server-side from the database; the storefront only says
	// which product, size and quantity.
	orderRepository := orders.NewPostgresRepository(databasePool)
	orderService := orders.NewService(orderRepository)
	orderHandler := orders.NewHandler(orderService)

	requestMux := http.NewServeMux()
	requestMux.Handle("GET /healthz", health.NewHandler(databasePool))
	requestMux.HandleFunc("GET /products", productHandler.List)
	requestMux.HandleFunc("GET /products/{slug}", productHandler.GetBySlug)
	requestMux.HandleFunc("POST /orders", orderHandler.Create)

	return &http.Server{
		Addr:              applicationConfig.APIAddress,
		Handler:           allowOrigin(applicationConfig.WebBaseURL, requestMux),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

// allowOrigin permits the configured storefront origin and handles CORS preflight requests.
func allowOrigin(allowedOrigin string, next http.Handler) http.Handler {
	// Apply CORS headers consistently before delegating actual requests.
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Origin") == allowedOrigin {
			responseWriter.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			responseWriter.Header().Set("Vary", "Origin")
		}

		if request.Method == http.MethodOptions {
			responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			responseWriter.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}
