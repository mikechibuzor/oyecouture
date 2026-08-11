package settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	keyAnnouncement   = "announcement"
	keyHomeContent    = "home_content"
	keyWhatsAppNumber = "whatsapp_number"

	maximumRequestBytes      = 32 << 10
	maximumAnnouncementChars = 180
	maximumHomeTextChars     = 320
	maximumHomePathChars     = 420
	maximumHomeTiles         = 6
)

var ErrInvalidSettings = errors.New("invalid settings")

type SiteSettings struct {
	Announcement   string `json:"announcement"`
	WhatsAppNumber string `json:"whatsappNumber"`
}

type HomeContent struct {
	Hero       HomeHeroContent     `json:"hero"`
	Categories HomeCategoryContent `json:"categories"`
}

type HomeHeroContent struct {
	HeadingLead string `json:"headingLead"`
	HeadingTail string `json:"headingTail"`
	Description string `json:"description"`
	ActionLabel string `json:"actionLabel"`
	ActionTo    string `json:"actionTo"`
}

type HomeCategoryContent struct {
	Eyebrow string             `json:"eyebrow"`
	Heading string             `json:"heading"`
	Tiles   []HomeCategoryTile `json:"tiles"`
}

type HomeCategoryTile struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Heading   string `json:"heading"`
	ImagePath string `json:"imagePath"`
	ImageAlt  string `json:"imageAlt"`
	To        string `json:"to"`
}

func DefaultHomeContent() HomeContent {
	return HomeContent{
		Hero: HomeHeroContent{
			HeadingLead: "Tailor Your Ambition",
			HeadingTail: "With Effortless Precision",
			Description: "Invest in pieces that reflect your ambition brilliantly. Bespoke tailoring that speaks before you do.",
			ActionLabel: "Discover Our Collection",
			ActionTo:    "/#collection",
		},
		Categories: HomeCategoryContent{
			Eyebrow: "Explore Oyecouture",
			Heading: "Shop By Category",
			Tiles: []HomeCategoryTile{
				{
					ID:        "suits",
					Label:     "Shop",
					Heading:   "Suits",
					ImagePath: "/images/catalogue/oye-cat-suits.webp",
					ImageAlt:  "Oyecouture tailored suit",
					To:        "/suits",
				},
				{
					ID:        "trads",
					Label:     "Shop",
					Heading:   "Trads",
					ImagePath: "https://res.cloudinary.com/exgithrw/image/upload/f_auto,q_auto,w_900/oyecouture/campaign-4",
					ImageAlt:  "Oyecouture traditional attire",
					To:        "/trad-attires",
				},
				{
					ID:        "casuals",
					Label:     "Shop",
					Heading:   "Casuals",
					ImagePath: "https://res.cloudinary.com/exgithrw/image/upload/f_auto,q_auto,w_900/oyecouture/casual-14",
					ImageAlt:  "Oyecouture casual clothing",
					To:        "/casuals/men",
				},
			},
		},
	}
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

func (repository *Repository) GetHomeContent(
	requestContext context.Context,
	defaults HomeContent,
) (HomeContent, error) {
	const query = `SELECT value FROM settings WHERE key = $1`

	var rawHomeContent string
	err := repository.pool.QueryRow(requestContext, query, keyHomeContent).Scan(&rawHomeContent)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return defaults, nil
		}

		return HomeContent{}, fmt.Errorf("query home content: %w", err)
	}

	var homeContent HomeContent
	if err := json.Unmarshal([]byte(rawHomeContent), &homeContent); err != nil {
		return HomeContent{}, fmt.Errorf("decode home content: %w", err)
	}

	normalizedHomeContent, err := normalizeHomeContent(homeContent)
	if err != nil {
		return HomeContent{}, fmt.Errorf("normalize stored home content: %w", err)
	}

	return normalizedHomeContent, nil
}

func (repository *Repository) UpdateHomeContent(
	requestContext context.Context,
	update HomeContent,
	defaults HomeContent,
) (HomeContent, error) {
	normalizedUpdate, err := normalizeHomeContent(update)
	if err != nil {
		return HomeContent{}, err
	}

	encodedUpdate, err := json.Marshal(normalizedUpdate)
	if err != nil {
		return HomeContent{}, fmt.Errorf("encode home content: %w", err)
	}

	const query = `
		INSERT INTO settings (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`

	if _, err := repository.pool.Exec(requestContext, query, keyHomeContent, string(encodedUpdate)); err != nil {
		return HomeContent{}, fmt.Errorf("update home content: %w", err)
	}

	return repository.GetHomeContent(requestContext, defaults)
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

func normalizeHomeContent(update HomeContent) (HomeContent, error) {
	normalizedUpdate := HomeContent{
		Hero: HomeHeroContent{
			HeadingLead: strings.TrimSpace(update.Hero.HeadingLead),
			HeadingTail: strings.TrimSpace(update.Hero.HeadingTail),
			Description: strings.TrimSpace(update.Hero.Description),
			ActionLabel: strings.TrimSpace(update.Hero.ActionLabel),
			ActionTo:    strings.TrimSpace(update.Hero.ActionTo),
		},
		Categories: HomeCategoryContent{
			Eyebrow: strings.TrimSpace(update.Categories.Eyebrow),
			Heading: strings.TrimSpace(update.Categories.Heading),
			Tiles:   make([]HomeCategoryTile, 0, len(update.Categories.Tiles)),
		},
	}

	requiredHeroFields := []string{
		normalizedUpdate.Hero.HeadingLead,
		normalizedUpdate.Hero.HeadingTail,
		normalizedUpdate.Hero.Description,
		normalizedUpdate.Hero.ActionLabel,
		normalizedUpdate.Hero.ActionTo,
		normalizedUpdate.Categories.Eyebrow,
		normalizedUpdate.Categories.Heading,
	}
	for _, value := range requiredHeroFields {
		if !validHomeText(value) {
			return HomeContent{}, ErrInvalidSettings
		}
	}
	if !validHomePath(normalizedUpdate.Hero.ActionTo) {
		return HomeContent{}, ErrInvalidSettings
	}
	if len(update.Categories.Tiles) == 0 || len(update.Categories.Tiles) > maximumHomeTiles {
		return HomeContent{}, ErrInvalidSettings
	}

	seenIDs := make(map[string]struct{}, len(update.Categories.Tiles))
	for _, tile := range update.Categories.Tiles {
		normalizedTile := HomeCategoryTile{
			ID:        strings.TrimSpace(tile.ID),
			Label:     strings.TrimSpace(tile.Label),
			Heading:   strings.TrimSpace(tile.Heading),
			ImagePath: strings.TrimSpace(tile.ImagePath),
			ImageAlt:  strings.TrimSpace(tile.ImageAlt),
			To:        strings.TrimSpace(tile.To),
		}

		if normalizedTile.ID == "" || !validHomeText(normalizedTile.Label) ||
			!validHomeText(normalizedTile.Heading) || !validHomeText(normalizedTile.ImageAlt) ||
			!validHomePath(normalizedTile.ImagePath) || !validHomePath(normalizedTile.To) {
			return HomeContent{}, ErrInvalidSettings
		}
		if _, exists := seenIDs[normalizedTile.ID]; exists {
			return HomeContent{}, ErrInvalidSettings
		}

		seenIDs[normalizedTile.ID] = struct{}{}
		normalizedUpdate.Categories.Tiles = append(normalizedUpdate.Categories.Tiles, normalizedTile)
	}

	return normalizedUpdate, nil
}

func validHomeText(value string) bool {
	return value != "" && len([]rune(value)) <= maximumHomeTextChars
}

func validHomePath(value string) bool {
	return value != "" && len([]rune(value)) <= maximumHomePathChars &&
		(strings.HasPrefix(value, "/") || strings.HasPrefix(value, "https://"))
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
	repository   *Repository
	defaults     SiteSettings
	homeDefaults HomeContent
}

func NewHandler(repository *Repository, defaults SiteSettings, homeDefaults HomeContent) *Handler {
	return &Handler{repository: repository, defaults: defaults, homeDefaults: homeDefaults}
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

func (handler *Handler) GetHomeContent(responseWriter http.ResponseWriter, request *http.Request) {
	homeContent, err := handler.repository.GetHomeContent(request.Context(), handler.homeDefaults)
	if err != nil {
		log.Printf("get home content: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error", "home content could not be loaded")
		return
	}

	writeJSON(responseWriter, http.StatusOK, homeContent)
}

func (handler *Handler) UpdateHomeContent(responseWriter http.ResponseWriter, request *http.Request) {
	var update HomeContent

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, maximumRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&update); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid_body", "the home content update could not be read")
		return
	}

	homeContent, err := handler.repository.UpdateHomeContent(request.Context(), update, handler.homeDefaults)
	if err != nil {
		if errors.Is(err, ErrInvalidSettings) {
			writeError(responseWriter, http.StatusBadRequest, "invalid_home_content", "please check the home content")
			return
		}

		log.Printf("update home content: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error", "home content could not be saved")
		return
	}

	writeJSON(responseWriter, http.StatusOK, homeContent)
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
