package settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	keyAnnouncement   = "announcement"
	keyWhatsAppNumber = "whatsapp_number"

	maximumRequestBytes      = 8 << 10
	maximumAnnouncementChars = 180
)

var ErrInvalidSettings = errors.New("invalid settings")

type SiteSettings struct {
	Announcement   string `json:"announcement"`
	WhatsAppNumber string `json:"whatsappNumber"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (repository *Repository) Get(
	requestContext context.Context,
	defaults SiteSettings,
) (SiteSettings, error) {
	const query = `
		SELECT key, value
		FROM settings
		WHERE key IN ($1, $2)`

	rows, err := repository.pool.Query(requestContext, query, keyAnnouncement, keyWhatsAppNumber)
	if err != nil {
		return SiteSettings{}, fmt.Errorf("query settings: %w", err)
	}
	defer rows.Close()

	currentSettings := defaults
	for rows.Next() {
		var key string
		var value string
		if err := rows.Scan(&key, &value); err != nil {
			return SiteSettings{}, fmt.Errorf("scan setting: %w", err)
		}

		switch key {
		case keyAnnouncement:
			currentSettings.Announcement = value
		case keyWhatsAppNumber:
			currentSettings.WhatsAppNumber = value
		}
	}

	return currentSettings, rows.Err()
}

func (repository *Repository) Update(
	requestContext context.Context,
	update SiteSettings,
	defaults SiteSettings,
) (SiteSettings, error) {
	normalizedUpdate, err := normalize(update)
	if err != nil {
		return SiteSettings{}, err
	}

	const query = `
		INSERT INTO settings (key, value)
		VALUES ($1, $2), ($3, $4)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`

	if _, err := repository.pool.Exec(
		requestContext,
		query,
		keyAnnouncement,
		normalizedUpdate.Announcement,
		keyWhatsAppNumber,
		normalizedUpdate.WhatsAppNumber,
	); err != nil {
		return SiteSettings{}, fmt.Errorf("update settings: %w", err)
	}

	return repository.Get(requestContext, defaults)
}

func normalize(update SiteSettings) (SiteSettings, error) {
	announcement := strings.TrimSpace(update.Announcement)
	whatsAppNumber := strings.TrimSpace(update.WhatsAppNumber)

	if len([]rune(announcement)) > maximumAnnouncementChars {
		return SiteSettings{}, ErrInvalidSettings
	}

	if whatsAppNumber != "" && len(digitsOnly(whatsAppNumber)) < 8 {
		return SiteSettings{}, ErrInvalidSettings
	}

	return SiteSettings{
		Announcement:   announcement,
		WhatsAppNumber: whatsAppNumber,
	}, nil
}

func digitsOnly(value string) string {
	var builder strings.Builder
	for _, character := range value {
		if character >= '0' && character <= '9' {
			builder.WriteRune(character)
		}
	}

	return builder.String()
}

type Handler struct {
	repository *Repository
	defaults   SiteSettings
}

func NewHandler(repository *Repository, defaults SiteSettings) *Handler {
	return &Handler{repository: repository, defaults: defaults}
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

func (handler *Handler) Get(responseWriter http.ResponseWriter, request *http.Request) {
	currentSettings, err := handler.repository.Get(request.Context(), handler.defaults)
	if err != nil {
		log.Printf("get settings: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error", "settings could not be loaded")
		return
	}

	writeJSON(responseWriter, http.StatusOK, currentSettings)
}

func (handler *Handler) Update(responseWriter http.ResponseWriter, request *http.Request) {
	var update SiteSettings

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, maximumRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&update); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid_body", "the settings update could not be read")
		return
	}

	currentSettings, err := handler.repository.Update(request.Context(), update, handler.defaults)
	if err != nil {
		if errors.Is(err, ErrInvalidSettings) {
			writeError(responseWriter, http.StatusBadRequest, "invalid_settings", "please check the settings")
			return
		}

		log.Printf("update settings: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error", "settings could not be saved")
		return
	}

	writeJSON(responseWriter, http.StatusOK, currentSettings)
}

func writeError(responseWriter http.ResponseWriter, statusCode int, errorCode string, errorMessage string) {
	writeJSON(responseWriter, statusCode, errorResponse{
		Error: errorDetail{Code: errorCode, Message: errorMessage},
	})
}

func writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		log.Printf("encode settings response: %v", err)
	}
}
