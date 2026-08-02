package orders

import (
	"context"
	"strings"
	"time"
)

// Admin-facing reads. These are separate from the storefront types on purpose:
// the owner sees contact details and payment state that a customer never should.

const (
	defaultPageSize   = 25
	maximumPageSize   = 100
	maximumSearchSize = 120
)

// Summary is one row of the admin orders list. It carries just enough to scan
// the list and decide what needs attention; the detail view loads the rest.
type Summary struct {
	OrderNo      string    `json:"orderNo"`
	Kind         Kind      `json:"kind"`
	Status       Status    `json:"status"`
	CustomerName string    `json:"customerName"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	ItemCount    int       `json:"itemCount"`
	TotalKobo    int64     `json:"totalKobo"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ListFilter narrows the admin orders list. The zero value lists everything,
// newest first.
type ListFilter struct {
	Status Status
	Kind   Kind
	Search string
	Limit  int
	Offset int
}

// Page is a slice of orders plus the total matching the filter, so the UI can
// show "showing 25 of 132" without a second request.
type Page struct {
	Orders []Summary `json:"orders"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

// validStatuses and validKinds mirror the database CHECK constraints. An
// unrecognised value is dropped rather than rejected: a stale bookmark should
// show all orders, not an error page.
var validStatuses = map[Status]struct{}{
	StatusPendingConfirmation: {}, StatusAwaitingPayment: {}, StatusPaid: {},
	StatusInProgress: {}, StatusShipped: {}, StatusDelivered: {},
	StatusCancelled: {}, StatusRefunded: {},
}

var validKinds = map[Kind]struct{}{KindBespoke: {}, KindReadyToWear: {}}

// AdminRepository is the persistence contract for admin order reads.
type AdminRepository interface {
	List(context.Context, ListFilter) (Page, error)
}

type AdminService struct {
	repository AdminRepository
}

func NewAdminService(repository AdminRepository) *AdminService {
	return &AdminService{repository: repository}
}

// List returns orders newest first, after clamping the filter to safe bounds.
func (service *AdminService) List(
	requestContext context.Context,
	filter ListFilter,
) (Page, error) {
	return service.repository.List(requestContext, normalizeFilter(filter))
}

// normalizeFilter drops unknown enum values and clamps paging, so a hand-typed
// query string can never ask for an unbounded read.
func normalizeFilter(filter ListFilter) ListFilter {
	if _, known := validStatuses[filter.Status]; !known {
		filter.Status = ""
	}
	if _, known := validKinds[filter.Kind]; !known {
		filter.Kind = ""
	}

	filter.Search = strings.TrimSpace(filter.Search)
	if len(filter.Search) > maximumSearchSize {
		filter.Search = filter.Search[:maximumSearchSize]
	}

	if filter.Limit <= 0 {
		filter.Limit = defaultPageSize
	}
	if filter.Limit > maximumPageSize {
		filter.Limit = maximumPageSize
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return filter
}
