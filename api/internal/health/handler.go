package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/database"
)

type Handler struct {
	database database.Pinger
}

type response struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func NewHandler(databasePinger database.Pinger) *Handler {
	return &Handler{database: databasePinger}
}

func (handler *Handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	requestContext, cancelRequest := context.WithTimeout(request.Context(), 2*time.Second)
	defer cancelRequest()

	healthResponse := response{
		Status: "ok",
		Services: map[string]string{
			"api":      "ok",
			"database": "ok",
		},
	}
	statusCode := http.StatusOK

	if err := handler.database.Ping(requestContext); err != nil {
		healthResponse.Status = "degraded"
		healthResponse.Services["database"] = "unavailable"
		statusCode = http.StatusServiceUnavailable
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	_ = json.NewEncoder(responseWriter).Encode(healthResponse)
}
