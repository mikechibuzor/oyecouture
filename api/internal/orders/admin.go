package orders

import (
	"context"
	"errors"
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

// Detail is the full admin read model for one order.
type Detail struct {
	OrderNo         string        `json:"orderNo"`
	Kind            Kind          `json:"kind"`
	Status          Status        `json:"status"`
	CustomerName    string        `json:"customerName"`
	Email           string        `json:"email"`
	Phone           string        `json:"phone"`
	Address         string        `json:"address,omitempty"`
	DeliveryZone    string        `json:"deliveryZone,omitempty"`
	DeliveryFeeKobo int64         `json:"deliveryFeeKobo"`
	SubtotalKobo    int64         `json:"subtotalKobo"`
	TotalKobo       int64         `json:"totalKobo"`
	PayToken        string        `json:"payToken"`
	PaymentProvider string        `json:"paymentProvider,omitempty"`
	PaymentRef      string        `json:"paymentRef,omitempty"`
	Items           []Item        `json:"items"`
	Measurements    *Measurements `json:"measurements,omitempty"`
	CreatedAt       time.Time     `json:"createdAt"`
	orderID         string
}

// StatusUpdate is the only mutable order field the first admin actions expose.
type StatusUpdate struct {
	Status Status `json:"status"`
}

// ListFilter narrows the admin orders list. The zero value lists everything,
// newest first. From and To are inclusive calendar days in the server's
// timezone, which is how the owner reads a date on screen.
type ListFilter struct {
	Status Status
	Kind   Kind
	Search string
	From   *time.Time
	To     *time.Time
	Limit  int
	Offset int
}

// DateFormat is the wire format for the From and To filters.
const DateFormat = "2006-01-02"

// ParseFilterDate reads a YYYY-MM-DD filter bound. An unparseable value is
// treated as absent rather than an error, so a hand-edited URL still lists.
func ParseFilterDate(value string) *time.Time {
	if value == "" {
		return nil
	}

	parsed, err := time.ParseInLocation(DateFormat, value, time.Local)
	if err != nil {
		return nil
	}

	return &parsed
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

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrInvalidStatusUpdate     = errors.New("order status is invalid")
	ErrInvalidStatusTransition = errors.New("order status transition is invalid")
	ErrPaymentMustBeServerSide = errors.New("paid status requires server-side payment verification")
)

// allowedAdminStatusTransitions encodes actions the owner can safely perform.
// `paid` is intentionally absent as a target: payment must come from the
// provider verification path, not from an admin click.
var allowedAdminStatusTransitions = map[Status]map[Status]struct{}{
	StatusPendingConfirmation: {
		StatusAwaitingPayment: {},
		StatusCancelled:       {},
	},
	StatusAwaitingPayment: {
		StatusCancelled: {},
	},
	StatusPaid: {
		StatusInProgress: {},
		StatusCancelled:  {},
		StatusRefunded:   {},
	},
	StatusInProgress: {
		StatusShipped:   {},
		StatusCancelled: {},
	},
	StatusShipped: {
		StatusDelivered: {},
	},
}

// Metrics is the summary strip above the orders list. Revenue counts only
// orders that have actually been paid for, so it never flatters the owner with
// money that has not arrived.
type Metrics struct {
	Total               int   `json:"total"`
	PendingConfirmation int   `json:"pendingConfirmation"`
	AwaitingPayment     int   `json:"awaitingPayment"`
	Paid                int   `json:"paid"`
	RevenueKobo         int64 `json:"revenueKobo"`
}

// PaidStatuses are the states an order can only reach after payment.
var PaidStatuses = []Status{StatusPaid, StatusInProgress, StatusShipped, StatusDelivered}

// MetricsPeriod is the window the summary strip covers.
type MetricsPeriod string

const (
	PeriodToday      MetricsPeriod = "today"
	PeriodLastWeek   MetricsPeriod = "last-week"
	PeriodLast14Days MetricsPeriod = "last-14-days"
	PeriodLast60Days MetricsPeriod = "last-60-days"
	PeriodLast90Days MetricsPeriod = "last-90-days"
	PeriodAllTime    MetricsPeriod = "all-time"
)

// periodDays maps each rolling window to its length. Today is absent because it
// runs from midnight rather than a fixed number of days back.
var periodDays = map[MetricsPeriod]int{
	PeriodLastWeek:   7,
	PeriodLast14Days: 14,
	PeriodLast60Days: 60,
	PeriodLast90Days: 90,
}

// AdminRepository is the persistence contract for admin order reads.
type AdminRepository interface {
	List(context.Context, ListFilter) (Page, error)
	Metrics(context.Context, *time.Time) (Metrics, error)
	Detail(context.Context, string) (Detail, error)
	UpdateStatus(context.Context, string, Status) error
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

// Metrics returns the counts shown above the orders list, over the requested
// window.
func (service *AdminService) Metrics(
	requestContext context.Context,
	period MetricsPeriod,
) (Metrics, error) {
	return service.repository.Metrics(requestContext, periodStart(period, time.Now()))
}

// Detail returns one order with its items and measurements.
func (service *AdminService) Detail(
	requestContext context.Context,
	orderNo string,
) (Detail, error) {
	return service.repository.Detail(requestContext, strings.TrimSpace(orderNo))
}

// UpdateStatus moves an order through owner-controlled lifecycle states.
func (service *AdminService) UpdateStatus(
	requestContext context.Context,
	orderNo string,
	update StatusUpdate,
) (Detail, error) {
	nextStatus := update.Status
	if nextStatus == StatusPaid {
		return Detail{}, ErrPaymentMustBeServerSide
	}
	if _, known := validStatuses[nextStatus]; !known {
		return Detail{}, ErrInvalidStatusUpdate
	}

	detail, err := service.Detail(requestContext, orderNo)
	if err != nil {
		return Detail{}, err
	}
	if !canAdminTransition(detail.Status, nextStatus) {
		return Detail{}, ErrInvalidStatusTransition
	}

	if err := service.repository.UpdateStatus(requestContext, detail.OrderNo, nextStatus); err != nil {
		return Detail{}, err
	}

	detail.Status = nextStatus
	return detail, nil
}

// canAdminTransition protects the paid invariant and blocks nonsensical jumps.
func canAdminTransition(currentStatus Status, nextStatus Status) bool {
	allowedTargets, hasTargets := allowedAdminStatusTransitions[currentStatus]
	if !hasTargets {
		return false
	}

	_, allowed := allowedTargets[nextStatus]
	return allowed
}

// periodStart resolves a window to the earliest order date it includes. A nil
// result means all time. An unrecognised period falls back to all time rather
// than erroring, so a stale bookmark still shows the summary.
func periodStart(period MetricsPeriod, now time.Time) *time.Time {
	if period == PeriodToday {
		// From local midnight, so "today" matches the owner's calendar day.
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return &midnight
	}

	days, isRollingWindow := periodDays[period]
	if !isRollingWindow {
		return nil
	}

	since := now.AddDate(0, 0, -days)

	return &since
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

	// "To" names a day the owner wants included, so the bound runs to the end of
	// it rather than to midnight at its start.
	if filter.To != nil {
		endOfDay := filter.To.AddDate(0, 0, 1)
		filter.To = &endOfDay
	}

	// A backwards range returns nothing and reads as a bug; drop the upper bound
	// instead so the list still shows something the owner can correct from.
	if filter.From != nil && filter.To != nil && filter.To.Before(*filter.From) {
		filter.To = nil
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
