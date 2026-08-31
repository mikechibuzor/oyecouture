package tracking

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/admin"
	"github.com/uzor-co/oyecouture/api/internal/ratelimit"
)

const (
	maximumTransitionBytes = 16 << 10
	publicTrackingWindow   = time.Hour
	publicTrackingLimit    = 120
)

type Handler struct {
	service *Service
	limiter *ratelimit.Limiter
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
		limiter: ratelimit.New(publicTrackingWindow, publicTrackingLimit),
	}
}

// Public serves a cache-resistant, PII-minimized order progress response.
func (handler *Handler) Public(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Cache-Control", "private, no-store")
	responseWriter.Header().Set("X-Robots-Tag", "noindex, nofollow")

	clientAddress := ratelimit.ClientAddress(request)
	if !handler.limiter.Allow(clientAddress) {
		writeError(responseWriter, http.StatusTooManyRequests, "too_many_requests",
			"too many tracking requests, try again shortly")
		return
	}
	handler.limiter.Record(clientAddress)

	state, err := handler.service.Public(request.Context(), request.PathValue("trackingToken"))
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, state)
}

// Admin returns the authenticated operational timeline for one order.
func (handler *Handler) Admin(responseWriter http.ResponseWriter, request *http.Request) {
	state, err := handler.service.Admin(request.Context(), request.PathValue("orderNo"))
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, state)
}

// Transition applies one optimistic, backend-approved fulfillment action.
func (handler *Handler) Transition(responseWriter http.ResponseWriter, request *http.Request) {
	session, authenticated := admin.SessionFrom(request.Context())
	if !authenticated {
		writeError(responseWriter, http.StatusUnauthorized, "not_authenticated", "not signed in")
		return
	}

	var input TransitionInput
	decoder := json.NewDecoder(http.MaxBytesReader(
		responseWriter, request.Body, maximumTransitionBytes,
	))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid_body",
			"the fulfillment update could not be read")
		return
	}

	state, err := handler.service.Transition(
		request.Context(), request.PathValue("orderNo"), session.User.ID, input,
	)
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, state)
}

func (handler *Handler) writeServiceError(responseWriter http.ResponseWriter, serviceError error) {
	switch {
	case errors.Is(serviceError, ErrNotFound):
		writeError(responseWriter, http.StatusNotFound, "tracking_not_found",
			"the order tracking link could not be found")
	case errors.Is(serviceError, ErrUnavailable):
		writeError(responseWriter, http.StatusConflict, "tracking_unavailable",
			"tracking becomes available after payment is confirmed")
	case errors.Is(serviceError, ErrInvalidTransition), errors.Is(serviceError, ErrInvalidInput):
		writeError(responseWriter, http.StatusBadRequest, "invalid_fulfillment_update",
			"that fulfillment update is not available")
	case errors.Is(serviceError, ErrVersionConflict):
		writeError(responseWriter, http.StatusConflict, "fulfillment_conflict",
			"the order changed in another session, reload and try again")
	default:
		log.Printf("tracking request failed: %v", serviceError)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error",
			"order tracking could not be loaded")
	}
}

func writeError(
	responseWriter http.ResponseWriter,
	statusCode int,
	code string,
	message string,
) {
	writeJSON(responseWriter, statusCode, map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}

func writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		log.Printf("encode tracking response: %v", err)
	}
}
