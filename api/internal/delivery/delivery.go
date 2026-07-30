package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Zone is a delivery destination and its fee. Fees live in the database so the
// owner can change them without a deploy.
type Zone struct {
	Name    string `json:"name"`
	FeeKobo int64  `json:"feeKobo"`
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
		// A missing row is a validation failure, not a server error.
		return 0, false, nil
	}

	return feeKobo, true, nil
}

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{repository: repository}
}

type listResponse struct {
	Zones []Zone `json:"zones"`
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
