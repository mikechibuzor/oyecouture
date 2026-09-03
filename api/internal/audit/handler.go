package audit

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) List(responseWriter http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	page, err := handler.service.List(request.Context(), limit, offset)
	if err != nil {
		log.Printf("list audit log: %v", err)
		writeJSON(responseWriter, http.StatusInternalServerError,
			map[string]string{"error": "the activity log could not be loaded"})
		return
	}

	writeJSON(responseWriter, http.StatusOK, page)
}

func writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		log.Printf("encode audit response: %v", err)
	}
}
