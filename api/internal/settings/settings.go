package settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	keyAnnouncement   = "announcement"
	keyHomeContent    = "home_content"
	keyWhatsAppNumber = "whatsapp_number"
	keySiteHeld       = "site_held"
	keyHoldHeadline   = "hold_headline"
	keyHoldMessage    = "hold_message"
	keyContactEmail   = "contact_email"
	keyNotifyEmail    = "notify_email"
	keyFromName       = "from_name"

	heldTrue = "true"

	maximumRequestBytes      = 32 << 10
	maximumAnnouncementChars = 180
	maximumHoldHeadlineChars = 70
	maximumHoldMessageChars  = 220
	maximumEmailChars        = 254
	maximumFromNameChars     = 60
	maximumHomeTextChars     = 320
	maximumHomePathChars     = 420
	maximumHomeTiles         = 6
)

var ErrInvalidSettings = errors.New("invalid settings")

// SiteSettings is what the storefront reads on every page, and therefore what any visitor
// can read. Everything here is meant to be published. The hold fields are public on purpose:
// while the shop is held the holding page is the only thing drawn, and it needs its own words.
type SiteSettings struct {
	Announcement   string `json:"announcement"`
	WhatsAppNumber string `json:"whatsappNumber"`
	// The address printed on the contact page for customers to write to.
	ContactEmail string `json:"contactEmail"`
	Held         bool   `json:"held"`
	HoldHeadline string `json:"holdHeadline"`
	HoldMessage  string `json:"holdMessage"`
}

// AdminSettings is the same settings plus the ones only the owner may see. It exists because
// the public and admin reads shared one payload: putting the notification inbox on
// SiteSettings would publish it to every visitor and every scraper.
type AdminSettings struct {
	SiteSettings
	// Where enquiries are delivered. Private: the customer writes to ContactEmail, which may
	// be a different address entirely.
	NotifyEmail string `json:"notifyEmail"`
	// The name recipients see in the From line of order, tracking and enquiry emails. The
	// address itself is not editable: it has to stay the one verified with the email provider.
	FromName string `json:"fromName"`
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

// Get reads every setting, including the owner-only ones. Callers serving the public decide
// what to publish by handing back only the embedded SiteSettings.
func (repository *Repository) Get(
	requestContext context.Context,
	defaults AdminSettings,
) (AdminSettings, error) {
	const query = `
		SELECT key, value
		FROM settings
		WHERE key IN ($1, $2, $3, $4, $5, $6, $7, $8)`

	rows, err := repository.pool.Query(
		requestContext,
		query,
		keyAnnouncement,
		keyWhatsAppNumber,
		keySiteHeld,
		keyHoldHeadline,
		keyHoldMessage,
		keyContactEmail,
		keyNotifyEmail,
		keyFromName,
	)
	if err != nil {
		return AdminSettings{}, fmt.Errorf("query settings: %w", err)
	}
	defer rows.Close()

	currentSettings := defaults
	for rows.Next() {
		var key string
		var value string
		if err := rows.Scan(&key, &value); err != nil {
			return AdminSettings{}, fmt.Errorf("scan setting: %w", err)
		}

		switch key {
		case keyAnnouncement:
			currentSettings.Announcement = value
		case keyWhatsAppNumber:
			currentSettings.WhatsAppNumber = value
		case keySiteHeld:
			currentSettings.Held = value == heldTrue
		case keyHoldHeadline:
			currentSettings.HoldHeadline = value
		case keyHoldMessage:
			currentSettings.HoldMessage = value
		case keyContactEmail:
			currentSettings.ContactEmail = value
		case keyNotifyEmail:
			currentSettings.NotifyEmail = value
		case keyFromName:
			currentSettings.FromName = value
		}
	}

	return currentSettings, rows.Err()
}

// NotifyEmail resolves where an enquiry should be delivered, reading only what the sender
// needs. The contact service calls it per send, so changing the address in the admin takes
// effect on the next enquiry rather than the next deploy.
func (repository *Repository) NotifyEmail(
	requestContext context.Context,
	fallback string,
) (string, error) {
	const query = `SELECT value FROM settings WHERE key = $1`

	var storedValue string
	err := repository.pool.QueryRow(requestContext, query, keyNotifyEmail).Scan(&storedValue)
	if errors.Is(err, pgx.ErrNoRows) {
		return fallback, nil
	}
	if err != nil {
		return "", fmt.Errorf("query notify email: %w", err)
	}

	if strings.TrimSpace(storedValue) == "" {
		return fallback, nil
	}

	return strings.TrimSpace(storedValue), nil
}

func (repository *Repository) Update(
	requestContext context.Context,
	update AdminSettings,
	defaults AdminSettings,
) (AdminSettings, error) {
	normalizedUpdate, err := normalize(update)
	if err != nil {
		return AdminSettings{}, err
	}

	const query = `
		INSERT INTO settings (key, value)
		VALUES ($1, $2), ($3, $4), ($5, $6), ($7, $8), ($9, $10), ($11, $12), ($13, $14), ($15, $16)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`

	if _, err := repository.pool.Exec(
		requestContext,
		query,
		keyAnnouncement,
		normalizedUpdate.Announcement,
		keyWhatsAppNumber,
		normalizedUpdate.WhatsAppNumber,
		keySiteHeld,
		boolAsSetting(normalizedUpdate.Held),
		keyHoldHeadline,
		normalizedUpdate.HoldHeadline,
		keyHoldMessage,
		normalizedUpdate.HoldMessage,
		keyContactEmail,
		normalizedUpdate.ContactEmail,
		keyNotifyEmail,
		normalizedUpdate.NotifyEmail,
		keyFromName,
		normalizedUpdate.FromName,
	); err != nil {
		return AdminSettings{}, fmt.Errorf("update settings: %w", err)
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

func normalize(update AdminSettings) (AdminSettings, error) {
	announcement := strings.TrimSpace(update.Announcement)
	whatsAppNumber := strings.TrimSpace(update.WhatsAppNumber)

	if len([]rune(announcement)) > maximumAnnouncementChars {
		return AdminSettings{}, ErrInvalidSettings
	}

	if whatsAppNumber != "" && len(digitsOnly(whatsAppNumber)) < 8 {
		return AdminSettings{}, ErrInvalidSettings
	}

	holdHeadline := strings.TrimSpace(update.HoldHeadline)
	holdMessage := strings.TrimSpace(update.HoldMessage)

	if len([]rune(holdHeadline)) > maximumHoldHeadlineChars ||
		len([]rune(holdMessage)) > maximumHoldMessageChars {
		return AdminSettings{}, ErrInvalidSettings
	}

	// A hold with no headline would draw an empty page, so the shop may only be held
	// once there is something to say on it.
	if update.Held && holdHeadline == "" {
		return AdminSettings{}, ErrInvalidSettings
	}

	contactEmail := strings.TrimSpace(update.ContactEmail)
	notifyEmail := strings.TrimSpace(update.NotifyEmail)

	// Either may be cleared, which falls back to the configured default rather than sending
	// nowhere. A value that is present must be a real address, or enquiries vanish silently.
	for _, address := range []string{contactEmail, notifyEmail} {
		if address == "" {
			continue
		}
		if len(address) > maximumEmailChars {
			return AdminSettings{}, ErrInvalidSettings
		}
		parsed, err := mail.ParseAddress(address)
		if err != nil || parsed.Address != address {
			return AdminSettings{}, ErrInvalidSettings
		}
	}

	fromName := strings.TrimSpace(update.FromName)
	if len([]rune(fromName)) > maximumFromNameChars {
		return AdminSettings{}, ErrInvalidSettings
	}

	return AdminSettings{
		SiteSettings: SiteSettings{
			Announcement:   announcement,
			WhatsAppNumber: whatsAppNumber,
			ContactEmail:   contactEmail,
			Held:           update.Held,
			HoldHeadline:   holdHeadline,
			HoldMessage:    holdMessage,
		},
		NotifyEmail: notifyEmail,
		FromName:    fromName,
	}, nil
}

// boolAsSetting stores the flag as the text this table holds for every other value.
func boolAsSetting(value bool) string {
	if value {
		return heldTrue
	}

	return "false"
}

// Held reports whether the shop is currently held, reading only the one flag the gate
// needs rather than the whole settings row.
func (repository *Repository) Held(requestContext context.Context) (bool, error) {
	const query = `SELECT value FROM settings WHERE key = $1`

	var storedValue string
	err := repository.pool.QueryRow(requestContext, query, keySiteHeld).Scan(&storedValue)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("query site hold: %w", err)
	}

	return storedValue == heldTrue, nil
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
	defaults     AdminSettings
	homeDefaults HomeContent
}

func NewHandler(repository *Repository, defaults AdminSettings, homeDefaults HomeContent) *Handler {
	return &Handler{repository: repository, defaults: defaults, homeDefaults: homeDefaults}
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

// Get serves the storefront. It returns only the embedded SiteSettings, so the owner's
// notification inbox is never part of a public response.
func (handler *Handler) Get(responseWriter http.ResponseWriter, request *http.Request) {
	currentSettings, err := handler.repository.Get(request.Context(), handler.defaults)
	if err != nil {
		log.Printf("get settings: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error", "settings could not be loaded")
		return
	}

	writeJSON(responseWriter, http.StatusOK, currentSettings.SiteSettings)
}

// GetForAdmin serves the settings screen, which may see everything. It is mounted behind
// ownerOnly; mounting it publicly would publish the notification inbox.
func (handler *Handler) GetForAdmin(responseWriter http.ResponseWriter, request *http.Request) {
	currentSettings, err := handler.repository.Get(request.Context(), handler.defaults)
	if err != nil {
		log.Printf("get settings: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error", "settings could not be loaded")
		return
	}

	writeJSON(responseWriter, http.StatusOK, currentSettings)
}

func (handler *Handler) Update(responseWriter http.ResponseWriter, request *http.Request) {
	var update AdminSettings

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

// FromAddress composes the From header for outgoing mail.
//
// Only the display name is editable. The address stays exactly as configured, because it is
// the one verified with the email provider and mail sent from anything else is rejected or
// lands in spam.
//
// Every failure returns the configured header unchanged. An order confirmation must not fail
// to send because a settings read did, so this can only ever improve the From line, never
// prevent one.
func (repository *Repository) FromAddress(
	requestContext context.Context,
	configured string,
) string {
	const query = `SELECT value FROM settings WHERE key = $1`

	var storedName string
	if err := repository.pool.QueryRow(requestContext, query, keyFromName).Scan(&storedName); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Printf("read from name: %v", err)
		}

		return configured
	}

	return ComposeFrom(strings.TrimSpace(storedName), configured)
}

// ComposeFrom puts a display name in front of the configured address, leaving the address
// alone. It is separate from the read so the rule can be tested without a database.
func ComposeFrom(displayName string, configured string) string {
	if displayName == "" {
		return configured
	}

	parsed, err := mail.ParseAddress(configured)
	if err != nil {
		return configured
	}

	return (&mail.Address{Name: displayName, Address: parsed.Address}).String()
}
