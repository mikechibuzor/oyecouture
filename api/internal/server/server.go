package server

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/uzor-co/oyecouture/api/internal/admin"
	"github.com/uzor-co/oyecouture/api/internal/audit"
	"github.com/uzor-co/oyecouture/api/internal/config"
	"github.com/uzor-co/oyecouture/api/internal/contact"
	"github.com/uzor-co/oyecouture/api/internal/delivery"
	emailpkg "github.com/uzor-co/oyecouture/api/internal/email"
	"github.com/uzor-co/oyecouture/api/internal/email/sendbyte"
	"github.com/uzor-co/oyecouture/api/internal/health"
	"github.com/uzor-co/oyecouture/api/internal/orders"
	"github.com/uzor-co/oyecouture/api/internal/payments"
	"github.com/uzor-co/oyecouture/api/internal/payments/paystack"
	"github.com/uzor-co/oyecouture/api/internal/products"
	"github.com/uzor-co/oyecouture/api/internal/settings"
	"github.com/uzor-co/oyecouture/api/internal/tracking"
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
	settingsRepository := settings.NewRepository(databasePool)
	settingsHandler := settings.NewHandler(settingsRepository, settings.SiteSettings{
		Announcement:   applicationConfig.Announcement,
		WhatsAppNumber: applicationConfig.WhatsAppNumber,
	}, settings.DefaultHomeContent())

	requestMux := http.NewServeMux()
	requestMux.Handle("GET /healthz", health.NewHandler(databasePool))
	requestMux.HandleFunc("GET /site-settings", settingsHandler.Get)
	requestMux.HandleFunc("GET /home-content", settingsHandler.GetHomeContent)
	requestMux.HandleFunc("GET /products", productHandler.List)
	requestMux.HandleFunc("GET /products/{slug}", productHandler.GetBySlug)
	requestMux.HandleFunc("POST /contact", contact.NewHandler(contact.NewService(
		emailProvider(applicationConfig),
		applicationConfig.EmailFrom,
		applicationConfig.OwnerNotifyEmail,
	)).Submit)
	requestMux.HandleFunc("POST /orders", orderHandler.Create)
	requestMux.HandleFunc("GET /delivery-zones", delivery.NewHandler(deliveryRepository).List)
	requestMux.HandleFunc("GET /payments/{payToken}", paymentHandler.Summary)
	requestMux.HandleFunc("POST /payments/{payToken}/initialize", paymentHandler.Initialize)
	requestMux.HandleFunc("POST /payments/{payToken}/verify", paymentHandler.Verify)
	requestMux.HandleFunc("POST /webhooks/paystack", paystackWebhookHandler.Handle)

	trackingService := tracking.NewService(
		tracking.NewPostgresRepository(databasePool),
		applicationConfig.EmailProvider,
		applicationConfig.WebBaseURL,
	)
	trackingHandler := tracking.NewHandler(trackingService)
	requestMux.HandleFunc("GET /tracking/{trackingToken}", trackingHandler.Public)

	// Ghost admin. Unlinked from the storefront and gated by a session cookie;
	// every route below RequireSession is unreachable without signing in.
	adminStore := admin.NewPostgresStore(databasePool)
	adminAuth := admin.NewAuthService(adminStore)
	adminHandler := admin.NewHandler(adminAuth, applicationConfig.AppEnv == "production")
	adminTeamHandler := admin.NewTeamHandler(admin.NewTeamService(adminStore))

	auditService := audit.NewService(audit.NewPostgresStore(databasePool))
	auditHandler := audit.NewHandler(auditService)

	// Every authenticated admin route goes through this rather than
	// RequireSession directly, so a route added later is audited without
	// anyone remembering to add it. RequireSession is the outer wrapper
	// because the recorder reads the session it puts in the context.
	authenticatedAdmin := func(next http.Handler) http.Handler {
		return adminHandler.RequireSession(auditService.Record(next))
	}

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
	requestMux.Handle("GET /admin/me", authenticatedAdmin(
		http.HandlerFunc(adminHandler.Me),
	))
	// Any signed-in admin may change their own password; it acts only on the
	// caller's own account. Everything else on the team desk is owner-only.
	requestMux.Handle("PATCH /admin/password", authenticatedAdmin(
		http.HandlerFunc(adminTeamHandler.ChangePassword),
	))
	requestMux.Handle("GET /admin/audit", authenticatedAdmin(
		admin.RequireOwner(http.HandlerFunc(auditHandler.List)),
	))
	requestMux.Handle("GET /admin/team", authenticatedAdmin(
		admin.RequireOwner(http.HandlerFunc(adminTeamHandler.List)),
	))
	requestMux.Handle("POST /admin/team", authenticatedAdmin(
		admin.RequireOwner(http.HandlerFunc(adminTeamHandler.Create)),
	))
	requestMux.Handle("PATCH /admin/team/{id}", authenticatedAdmin(
		admin.RequireOwner(http.HandlerFunc(adminTeamHandler.Update)),
	))
	requestMux.Handle("DELETE /admin/team/{id}", authenticatedAdmin(
		admin.RequireOwner(http.HandlerFunc(adminTeamHandler.Delete)),
	))
	requestMux.Handle("GET /admin/settings", authenticatedAdmin(
		http.HandlerFunc(settingsHandler.Get),
	))
	requestMux.Handle("PATCH /admin/settings", authenticatedAdmin(
		http.HandlerFunc(settingsHandler.Update),
	))
	requestMux.Handle("GET /admin/home-content", authenticatedAdmin(
		http.HandlerFunc(settingsHandler.GetHomeContent),
	))
	requestMux.Handle("PATCH /admin/home-content", authenticatedAdmin(
		http.HandlerFunc(settingsHandler.UpdateHomeContent),
	))

	adminOrderHandler := orders.NewAdminHandler(orders.NewAdminService(orderRepository))
	imageUploader := products.NewCloudinaryUploader(
		applicationConfig.CloudinaryCloudName,
		applicationConfig.CloudinaryAPIKey,
		applicationConfig.CloudinaryAPISecret,
		applicationConfig.CloudinaryFolder,
	)
	adminProductHandler := products.NewAdminHandler(
		products.NewAdminServiceWithUploader(
			productRepository,
			imageUploader,
			applicationConfig.R2PublicBase,
		),
	)
	requestMux.Handle("GET /admin/orders", authenticatedAdmin(
		http.HandlerFunc(adminOrderHandler.List),
	))
	requestMux.Handle("GET /admin/orders/metrics", authenticatedAdmin(
		http.HandlerFunc(adminOrderHandler.Metrics),
	))
	requestMux.Handle("GET /admin/orders/{orderNo}", authenticatedAdmin(
		http.HandlerFunc(adminOrderHandler.Detail),
	))
	requestMux.Handle("PATCH /admin/orders/{orderNo}/status", authenticatedAdmin(
		http.HandlerFunc(adminOrderHandler.UpdateStatus),
	))
	requestMux.Handle("GET /admin/orders/{orderNo}/tracking", authenticatedAdmin(
		http.HandlerFunc(trackingHandler.Admin),
	))
	requestMux.Handle("PATCH /admin/orders/{orderNo}/tracking", authenticatedAdmin(
		http.HandlerFunc(trackingHandler.Transition),
	))
	requestMux.Handle("GET /admin/products", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.List),
	))
	requestMux.Handle("POST /admin/products", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.Create),
	))
	requestMux.Handle("GET /admin/products/{slug}", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.GetBySlug),
	))
	requestMux.Handle("PUT /admin/products/{slug}", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.Update),
	))
	requestMux.Handle("DELETE /admin/products/{slug}", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.Delete),
	))
	requestMux.Handle("PATCH /admin/products/{slug}/active", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.UpdateActive),
	))
	adminDeliveryHandler := delivery.NewAdminHandler(deliveryRepository)
	requestMux.Handle("GET /admin/delivery-zones", authenticatedAdmin(
		http.HandlerFunc(adminDeliveryHandler.List),
	))
	requestMux.Handle("POST /admin/delivery-zones", authenticatedAdmin(
		http.HandlerFunc(adminDeliveryHandler.Create),
	))
	requestMux.Handle("PATCH /admin/delivery-zones/{name}", authenticatedAdmin(
		http.HandlerFunc(adminDeliveryHandler.Update),
	))
	requestMux.Handle("GET /admin/products/{slug}/images", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.Images),
	))
	requestMux.Handle("POST /admin/products/{slug}/images", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.UploadImage),
	))
	requestMux.Handle("PATCH /admin/products/{slug}/images/reorder", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.ReorderImages),
	))
	requestMux.Handle("DELETE /admin/products/{slug}/images/{imageID}", authenticatedAdmin(
		http.HandlerFunc(adminProductHandler.DeleteImage),
	))

	httpServer := &http.Server{
		Addr:              applicationConfig.APIAddress,
		Handler:           allowOrigin(applicationConfig.WebBaseURL, applicationConfig.AppEnv, requestMux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	if emailService != nil {
		workerContext, stopWorker := context.WithCancel(context.Background())
		httpServer.RegisterOnShutdown(stopWorker)
		go emailService.RunFulfillmentWorker(workerContext)
	}

	return httpServer
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
) *emailpkg.Service {
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
func allowOrigin(allowedOrigin string, appEnv string, next http.Handler) http.Handler {
	// Apply CORS headers consistently before delegating actual requests.
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		requestOrigin := request.Header.Get("Origin")
		if isAllowedOrigin(requestOrigin, allowedOrigin, appEnv) {
			responseWriter.Header().Set("Access-Control-Allow-Origin", requestOrigin)
			responseWriter.Header().Set("Vary", "Origin")
			// The admin session cookie only travels on credentialed requests.
			// Development runs the storefront and API on different ports; in
			// production Caddy serves both from one origin and CORS never applies.
			responseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if request.Method == http.MethodOptions {
			responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			responseWriter.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}

func isAllowedOrigin(requestOrigin string, allowedOrigin string, appEnv string) bool {
	if requestOrigin == "" {
		return false
	}
	if requestOrigin == allowedOrigin {
		return true
	}
	if appEnv == "production" {
		return false
	}

	parsedOrigin, err := url.Parse(requestOrigin)
	if err != nil {
		return false
	}

	return parsedOrigin.Scheme == "http" &&
		(parsedOrigin.Hostname() == "localhost" || parsedOrigin.Hostname() == "127.0.0.1")
}
