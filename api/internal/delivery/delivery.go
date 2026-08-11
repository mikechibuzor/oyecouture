package delivery

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

var (
	ErrInvalidZoneUpdate = errors.New("invalid delivery zone update")
	ErrZoneNotFound      = errors.New("delivery zone not found")
)

// Zone is a delivery destination and its fee. Fees live in the database so the
// owner can change them without a deploy.
type Zone struct {
	Name    string `json:"name"`
	FeeKobo int64  `json:"feeKobo"`
}

type AdminZone struct {
	Name    string `json:"name"`
	FeeKobo int64  `json:"feeKobo"`
	Active  bool   `json:"active"`
}

type AdminZoneUpdate struct {
	FeeKobo int64 `json:"feeKobo"`
	Active  *bool `json:"active"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// ActiveZones returns every selectable zone, cheapest first.
func (repository *Repository) ActiveZones(requestContext context.Context) ([]Zone, error) {
	const query = `SELECT name, fee_kobo FROM delivery_zones WHERE active = TRUE ORDER BY fee_kobo`

	rows, err := repository.pool.Query(requestContext, query)
	if err != nil {
		return nil, fmt.Errorf("query delivery zones: %w", err)
	}
	defer rows.Close()

	zones := make([]Zone, 0, 3)
	for rows.Next() {
		var zone Zone
		if err := rows.Scan(&zone.Name, &zone.FeeKobo); err != nil {
			return nil, fmt.Errorf("scan delivery zone: %w", err)
		}
		zones = append(zones, zone)
	}

	return zones, rows.Err()
}

// FeeForZone returns the fee for one zone, and whether that zone exists. Orders
// use this so the delivery fee is never taken from the client.
func (repository *Repository) FeeForZone(
	requestContext context.Context,
	zoneName string,
) (int64, bool, error) {
	const query = `SELECT fee_kobo FROM delivery_zones WHERE active = TRUE AND name = $1`

	var feeKobo int64
	err := repository.pool.QueryRow(requestContext, query, zoneName).Scan(&feeKobo)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, false, fmt.Errorf("query delivery zone fee: %w", err)
		}

		// A missing row is a validation failure, not a server error.
		return 0, false, nil
	}

	return feeKobo, true, nil
}

func (repository *Repository) AdminZones(requestContext context.Context) ([]AdminZone, error) {
	const query = `
		SELECT name, fee_kobo, active
		FROM delivery_zones
		ORDER BY active DESC, fee_kobo, name`

	rows, err := repository.pool.Query(requestContext, query)
	if err != nil {
		return nil, fmt.Errorf("query admin delivery zones: %w", err)
	}
	defer rows.Close()

	zones := make([]AdminZone, 0, 3)
	for rows.Next() {
		var zone AdminZone
		if err := rows.Scan(&zone.Name, &zone.FeeKobo, &zone.Active); err != nil {
			return nil, fmt.Errorf("scan admin delivery zone: %w", err)
		}
		zones = append(zones, zone)
	}

	return zones, rows.Err()
}

func (repository *Repository) UpdateAdminZone(
	requestContext context.Context,
	zoneName string,
	update AdminZoneUpdate,
) (AdminZone, error) {
	cleanZoneName := strings.TrimSpace(zoneName)
	if cleanZoneName == "" || update.FeeKobo < 0 || update.Active == nil {
		return AdminZone{}, ErrInvalidZoneUpdate
	}

	commandTag, err := repository.pool.Exec(
		requestContext,
		`UPDATE delivery_zones SET fee_kobo = $2, active = $3 WHERE name = $1`,
		cleanZoneName,
		update.FeeKobo,
		*update.Active,
	)
	if err != nil {
		return AdminZone{}, fmt.Errorf("update admin delivery zone: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return AdminZone{}, ErrZoneNotFound
	}

	return repository.AdminZoneByName(requestContext, cleanZoneName)
}

func (repository *Repository) AdminZoneByName(
	requestContext context.Context,
	zoneName string,
) (AdminZone, error) {
	const query = `SELECT name, fee_kobo, active FROM delivery_zones WHERE name = $1`

	var zone AdminZone
	err := repository.pool.QueryRow(requestContext, query, zoneName).Scan(
		&zone.Name,
		&zone.FeeKobo,
		&zone.Active,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AdminZone{}, ErrZoneNotFound
		}
		return AdminZone{}, fmt.Errorf("query admin delivery zone: %w", err)
	}

	return zone, nil
}

type Handler struct {
	repository *Repository
}

type AdminHandler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{repository: repository}
}

func NewAdminHandler(repository *Repository) *AdminHandler {
	return &AdminHandler{repository: repository}
}

type listResponse struct {
	Zones []Zone `json:"zones"`
}

type adminListResponse struct {
	Zones []AdminZone `json:"zones"`
}

type adminDetailResponse struct {
	Zone AdminZone `json:"zone"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

// List exposes the selectable delivery zones to the checkout page.
func (handler *Handler) List(responseWriter http.ResponseWriter, request *http.Request) {
	zones, err := handler.repository.ActiveZones(request.Context())
	if err != nil {
		log.Printf("list delivery zones: %v", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(responseWriter).Encode(listResponse{Zones: zones}); err != nil {
		log.Printf("encode delivery zones: %v", err)
	}
}

func (handler *AdminHandler) List(responseWriter http.ResponseWriter, request *http.Request) {
	zones, err := handler.repository.AdminZones(request.Context())
	if err != nil {
		log.Printf("list admin delivery zones: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error", "delivery zones could not be loaded")
		return
	}
	if zones == nil {
		zones = make([]AdminZone, 0)
	}

	writeJSON(responseWriter, http.StatusOK, adminListResponse{Zones: zones})
}

func (handler *AdminHandler) Update(responseWriter http.ResponseWriter, request *http.Request) {
	var update AdminZoneUpdate

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, 8<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&update); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid_body", "the delivery update could not be read")
		return
	}

	zone, err := handler.repository.UpdateAdminZone(
		request.Context(),
		request.PathValue("name"),
		update,
	)
	if err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, adminDetailResponse{Zone: zone})
}

func (handler *AdminHandler) writeServiceError(responseWriter http.ResponseWriter, serviceError error) {
	switch {
	case errors.Is(serviceError, ErrInvalidZoneUpdate):
		writeError(responseWriter, http.StatusBadRequest, "invalid_delivery_zone", "please check the delivery fee")
	case errors.Is(serviceError, ErrZoneNotFound):
		writeError(responseWriter, http.StatusNotFound, "delivery_zone_not_found", "delivery zone was not found")
	default:
		log.Printf("update admin delivery zone: %v", serviceError)
		writeError(responseWriter, http.StatusInternalServerError, "internal_error", "delivery zone could not be updated")
	}
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
		log.Printf("encode delivery response: %v", err)
	}
}
