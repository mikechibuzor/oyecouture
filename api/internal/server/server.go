package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/uzor-co/oyecouture/api/internal/admin"
	"github.com/uzor-co/oyecouture/api/internal/config"
	"github.com/uzor-co/oyecouture/api/internal/delivery"
	emailpkg "github.com/uzor-co/oyecouture/api/internal/email"
	"github.com/uzor-co/oyecouture/api/internal/email/sendbyte"
	"github.com/uzor-co/oyecouture/api/internal/health"
	"github.com/uzor-co/oyecouture/api/internal/orders"
	"github.com/uzor-co/oyecouture/api/internal/payments"
	"github.com/uzor-co/oyecouture/api/internal/payments/paystack"
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
	deliveryRepository := delivery.NewRepository(databasePool)
	orderRepository := orders.NewPostgresRepository(databasePool)
	orderService := orders.NewService(orderRepository, deliveryRepository)
	orderHandler := orders.NewHandler(orderService)

	paymentRepository := payments.NewPostgresRepository(databasePool)
	paymentProvider := paymentProvider(applicationConfig)
	emailService := emailService(applicationConfig, databasePool)
	paymentService := payments.NewService(
		paymentRepository, paymentProvider, applicationConfig.WebBaseURL, emailService,
	)
	paymentHandler := payments.NewHandler(paymentService)
	paystackWebhookHandler := paystack.NewWebhookHandler(
		paymentService, applicationConfig.PaystackSecretKey,
	)

	requestMux := http.NewServeMux()
	requestMux.Handle("GET /healthz", health.NewHandler(databasePool))
	requestMux.HandleFunc("GET /products", productHandler.List)
	requestMux.HandleFunc("GET /products/{slug}", productHandler.GetBySlug)
	requestMux.HandleFunc("POST /orders", orderHandler.Create)
	requestMux.HandleFunc("GET /delivery-zones", delivery.NewHandler(deliveryRepository).List)
	requestMux.HandleFunc("GET /payments/{payToken}", paymentHandler.Summary)
	requestMux.HandleFunc("POST /payments/{payToken}/initialize", paymentHandler.Initialize)
	requestMux.HandleFunc("POST /payments/{payToken}/verify", paymentHandler.Verify)
	requestMux.HandleFunc("POST /webhooks/paystack", paystackWebhookHandler.Handle)

	// Ghost admin. Unlinked from the storefront and gated by a session cookie;
	// every route below RequireSession is unreachable without signing in.
	adminStore := admin.NewPostgresStore(databasePool)
	adminAuth := admin.NewAuthService(adminStore)
	adminHandler := admin.NewHandler(adminAuth, applicationConfig.AppEnv == "production")

	// Clear dead sessions once at startup rather than running a scheduler.
	go func() {
		cleanupContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := adminStore.DeleteExpiredSessions(cleanupContext); err != nil {
			log.Printf("clear expired admin sessions: %v", err)
		}
	}()

	requestMux.HandleFunc("POST /admin/login", adminHandler.LogIn)
	requestMux.HandleFunc("POST /admin/logout", adminHandler.LogOut)
	requestMux.Handle("GET /admin/me", adminHandler.RequireSession(
		http.HandlerFunc(adminHandler.Me),
	))

	adminOrderHandler := orders.NewAdminHandler(orders.NewAdminService(orderRepository))
	requestMux.Handle("GET /admin/orders", adminHandler.RequireSession(
		http.HandlerFunc(adminOrderHandler.List),
	))
	requestMux.Handle("GET /admin/orders/metrics", adminHandler.RequireSession(
		http.HandlerFunc(adminOrderHandler.Metrics),
	))
	requestMux.Handle("GET /admin/orders/{orderNo}", adminHandler.RequireSession(
		http.HandlerFunc(adminOrderHandler.Detail),
	))
	requestMux.Handle("PATCH /admin/orders/{orderNo}/status", adminHandler.RequireSession(
		http.HandlerFunc(adminOrderHandler.UpdateStatus),
	))

	return &http.Server{
		Addr:              applicationConfig.APIAddress,
		Handler:           allowOrigin(applicationConfig.WebBaseURL, requestMux),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

// paymentProvider selects the configured gateway while keeping unsupported
// providers unavailable rather than silently charging through the wrong one.
func paymentProvider(applicationConfig config.Config) payments.Provider {
	if applicationConfig.PaymentProvider == payments.ProviderPaystack {
		return paystack.NewPaystack(applicationConfig.PaystackSecretKey)
	}

	return nil
}

// emailService builds the paid-order email sender only when email config is complete.
func emailService(
	applicationConfig config.Config,
	databasePool *pgxpool.Pool,
) payments.PaidOrderEmailSender {
	emailProvider := emailProvider(applicationConfig)
	if emailProvider == nil || applicationConfig.EmailFrom == "" {
		return nil
	}

	return emailpkg.NewService(
		emailpkg.NewPostgresRepository(databasePool),
		emailProvider,
		applicationConfig.EmailFrom,
		applicationConfig.WebBaseURL,
	)
}

// emailProvider selects the configured transactional email provider.
func emailProvider(applicationConfig config.Config) emailpkg.Provider {
	if applicationConfig.EmailProvider == emailpkg.ProviderSendByte &&
		applicationConfig.SendByteAPIKey != "" {
		return sendbyte.NewSendByte(applicationConfig.SendByteAPIKey)
	}

	return nil
}

// allowOrigin permits the configured storefront origin and handles CORS preflight requests.
func allowOrigin(allowedOrigin string, next http.Handler) http.Handler {
	// Apply CORS headers consistently before delegating actual requests.
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Origin") == allowedOrigin {
			responseWriter.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			responseWriter.Header().Set("Vary", "Origin")
			// The admin session cookie only travels on credentialed requests.
			// Development runs the storefront and API on different ports; in
			// production Caddy serves both from one origin and CORS never applies.
			responseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if request.Method == http.MethodOptions {
			responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			responseWriter.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}
