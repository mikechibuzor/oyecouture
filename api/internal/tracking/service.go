package tracking

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/email"
)

const (
	trackingTokenBytes        = 32
	maximumPublicNoteLength   = 500
	maximumCarrierFieldLength = 200
	dateFormat                = "2006-01-02"
	// Mirrors the CHECK constraints in migration 000017, so an over-long value
	// is refused with a clear error rather than a database violation.
	maximumSignedForByLength      = 120
	maximumDeliveryLocationLength = 240
	// A handover happens at a time of day, not just on a date, so the delivered
	// panel needs minutes. Accepts what a datetime-local input submits.
	dateTimeFormat = "2006-01-02T15:04"
)

var (
	ErrNotFound          = errors.New("tracking not found")
	ErrUnavailable       = errors.New("order tracking is not available")
	ErrInvalidTransition = errors.New("fulfillment transition is invalid")
	ErrVersionConflict   = errors.New("fulfillment version conflict")
	ErrInvalidInput      = errors.New("fulfillment input is invalid")
)

type Repository interface {
	Admin(context.Context, string) (AdminState, error)
	Public(context.Context, string) (PublicState, error)
	Transition(context.Context, TransitionCommand) error
}

// TransitionCommand is the validated, transaction-ready fulfillment mutation.
type TransitionCommand struct {
	OrderNo               string
	CurrentStage          Stage
	NextStage             Stage
	ExpectedVersion       int
	PublicNote            string
	EstimatedDeliveryDate *time.Time
	CarrierName           string
	CarrierTrackingNumber string
	CarrierTrackingURL    string
	SignedForBy           string
	SignedAt              *time.Time
	DeliveryLocation      string
	TrackingToken         string
	ActorAdminID          string
	LegacyOrderStatus     string
	OccurredAt            time.Time
	Notification          *email.Event
}

type Service struct {
	repository        Repository
	emailProviderName string
	webBaseURL        string
}

func NewService(repository Repository, emailProviderName string, webBaseURL string) *Service {
	return &Service{
		repository:        repository,
		emailProviderName: strings.TrimSpace(emailProviderName),
		webBaseURL:        strings.TrimRight(webBaseURL, "/"),
	}
}

// Admin returns the full audited timeline and the next legal actions.
func (service *Service) Admin(
	requestContext context.Context,
	orderNo string,
) (AdminState, error) {
	state, err := service.repository.Admin(requestContext, strings.TrimSpace(orderNo))
	if err != nil {
		return AdminState{}, err
	}
	if state.CurrentStage == "" {
		return AdminState{}, ErrUnavailable
	}

	return service.presentAdminState(state), nil
}

// Public resolves a bearer tracking link into a PII-minimized timeline.
func (service *Service) Public(
	requestContext context.Context,
	trackingToken string,
) (PublicState, error) {
	trackingToken = strings.TrimSpace(trackingToken)
	if len(trackingToken) < 40 {
		return PublicState{}, ErrNotFound
	}

	state, err := service.repository.Public(requestContext, trackingToken)
	if err != nil {
		return PublicState{}, err
	}
	state.CurrentStageLabel = StageLabel(state.CurrentStage)
	for eventIndex := range state.Events {
		state.Events[eventIndex].Label = StageLabel(state.Events[eventIndex].Stage)
		state.Events[eventIndex].ActorEmail = ""
	}

	return state, nil
}

// Transition advances an order once, records the event, and enqueues any
// milestone email in the same database transaction.
func (service *Service) Transition(
	requestContext context.Context,
	orderNo string,
	actorAdminID string,
	input TransitionInput,
) (AdminState, error) {
	state, err := service.repository.Admin(requestContext, strings.TrimSpace(orderNo))
	if err != nil {
		return AdminState{}, err
	}
	if input.ExpectedVersion != state.Version {
		return AdminState{}, ErrVersionConflict
	}
	if !containsStage(nextStages(state.Kind, state.CurrentStage), input.Stage) {
		return AdminState{}, ErrInvalidTransition
	}

	publicNote := strings.TrimSpace(input.PublicNote)
	if len(publicNote) > maximumPublicNoteLength {
		return AdminState{}, ErrInvalidInput
	}
	estimatedDeliveryDate, err := parseOptionalDate(input.EstimatedDeliveryDate)
	if err != nil {
		return AdminState{}, ErrInvalidInput
	}
	carrierName := strings.TrimSpace(input.CarrierName)
	carrierTrackingNumber := strings.TrimSpace(input.CarrierTrackingNumber)
	carrierTrackingURL := strings.TrimSpace(input.CarrierTrackingURL)
	if !validCarrierFields(carrierName, carrierTrackingNumber, carrierTrackingURL) {
		return AdminState{}, ErrInvalidInput
	}

	// Handover details belong to delivery and nothing else. Accepting them on
	// any other transition would let an order carry a signature before it had
	// been handed to anyone.
	signedForBy := strings.TrimSpace(input.SignedForBy)
	deliveryLocation := strings.TrimSpace(input.DeliveryLocation)
	signedAt, err := parseOptionalTime(input.SignedAt)
	if err != nil {
		return AdminState{}, ErrInvalidInput
	}
	if input.Stage != StageDelivered &&
		(signedForBy != "" || deliveryLocation != "" || signedAt != nil) {
		return AdminState{}, ErrInvalidInput
	}
	if len(signedForBy) > maximumSignedForByLength ||
		len(deliveryLocation) > maximumDeliveryLocationLength {
		return AdminState{}, ErrInvalidInput
	}

	trackingToken := state.trackingToken
	if input.Stage == StageDispatched && trackingToken == "" {
		trackingToken, err = newTrackingToken()
		if err != nil {
			return AdminState{}, fmt.Errorf("generate tracking token: %w", err)
		}
	}

	occurredAt := time.Now().UTC()
	command := TransitionCommand{
		OrderNo:               state.OrderNo,
		CurrentStage:          state.CurrentStage,
		NextStage:             input.Stage,
		ExpectedVersion:       input.ExpectedVersion,
		PublicNote:            publicNote,
		EstimatedDeliveryDate: estimatedDeliveryDate,
		SignedForBy:           signedForBy,
		SignedAt:              signedAt,
		DeliveryLocation:      deliveryLocation,
		CarrierName:           carrierName,
		CarrierTrackingNumber: carrierTrackingNumber,
		CarrierTrackingURL:    carrierTrackingURL,
		TrackingToken:         trackingToken,
		ActorAdminID:          actorAdminID,
		LegacyOrderStatus:     legacyOrderStatus(input.Stage),
		OccurredAt:            occurredAt,
	}
	command.Notification, err = service.notification(state, command)
	if err != nil {
		return AdminState{}, err
	}

	if err := service.repository.Transition(requestContext, command); err != nil {
		return AdminState{}, err
	}

	return service.Admin(requestContext, state.OrderNo)
}

// emailItems converts tracking items for the email package, carrying the line
// totals the dispatch and delivery designs show.
func emailItems(items []Item) []email.Item {
	converted := make([]email.Item, 0, len(items))
	for _, item := range items {
		converted = append(converted, email.Item{
			Name:     item.Name,
			Size:     item.Size,
			Quantity: item.Quantity,
			LineKobo: item.LineKobo,
		})
	}

	return converted
}

// notification builds only the milestone emails approved for the MVP.
func (service *Service) notification(
	state AdminState,
	command TransitionCommand,
) (*email.Event, error) {
	if service.emailProviderName == "" || state.customerEmail == "" {
		return nil, nil
	}

	emailType := ""
	subject := ""
	switch command.NextStage {
	case StageInProduction:
		emailType = email.TypeFulfillmentUpdate
		subject = "Your Oyecouture piece is in production - " + state.OrderNo
	case StageDispatched:
		emailType = email.TypeOrderDispatched
		subject = "Your Oyecouture order is on its way - " + state.OrderNo
	case StageDelivered:
		emailType = email.TypeOrderDelivered
		subject = "Your Oyecouture order has arrived - " + state.OrderNo
	default:
		return nil, nil
	}

	payload := email.FulfillmentPayload{
		CustomerName:          state.customerName,
		OrderNo:               state.OrderNo,
		Stage:                 string(command.NextStage),
		StageLabel:            StageLabel(command.NextStage),
		PublicNote:            command.PublicNote,
		CarrierName:           command.CarrierName,
		CarrierTrackingNumber: command.CarrierTrackingNumber,
		CarrierTrackingURL:    command.CarrierTrackingURL,
		OccurredAt:            command.OccurredAt,
		SubtotalKobo:          state.subtotalKobo,
		DeliveryFeeKobo:       state.deliveryFeeKobo,
		TotalKobo:             state.totalKobo,
		Items:                 emailItems(state.items),
		SignedForBy:           command.SignedForBy,
		DeliveryLocation:      command.DeliveryLocation,
	}
	if command.SignedAt != nil {
		payload.SignedAt = command.SignedAt.Format("2 January 2006, 15:04")
	}
	if command.EstimatedDeliveryDate != nil {
		payload.EstimatedDeliveryDate = command.EstimatedDeliveryDate.Format("2 Jan 2006")
	}
	if command.TrackingToken != "" &&
		(command.NextStage == StageDispatched || command.NextStage == StageDelivered) {
		payload.TrackingURL = service.trackingURL(command.TrackingToken)
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode fulfillment notification: %w", err)
	}

	return &email.Event{
		OrderID:  state.orderID,
		Type:     emailType,
		Provider: service.emailProviderName,
		IdempotencyKey: fmt.Sprintf(
			"order:%s:fulfillment:%d:%s", state.OrderNo, state.Version+1, command.NextStage,
		),
		RecipientEmail: state.customerEmail,
		Subject:        subject,
		Payload:        encodedPayload,
	}, nil
}

func (service *Service) presentAdminState(state AdminState) AdminState {
	state.CurrentStageLabel = StageLabel(state.CurrentStage)
	state.AllowedNextStages = nextStages(state.Kind, state.CurrentStage)
	if state.trackingToken != "" {
		state.TrackingURL = service.trackingURL(state.trackingToken)
	}
	for eventIndex := range state.Events {
		state.Events[eventIndex].Label = StageLabel(state.Events[eventIndex].Stage)
	}

	return state
}

func (service *Service) trackingURL(token string) string {
	return service.webBaseURL + "/track/" + url.PathEscape(token)
}

// StageLabel is shared by admin and customer responses so wording cannot drift.
func StageLabel(stage Stage) string {
	switch stage {
	case StagePaymentConfirmed:
		return "Payment confirmed"
	case StageMeasurementsApproved:
		return "Measurements approved"
	case StagePreparingOrder:
		return "Preparing order"
	case StageInProduction:
		return "In production"
	case StageQualityCheck:
		return "Quality check"
	case StageReadyForDispatch:
		return "Ready for dispatch"
	case StageDispatched:
		return "Dispatched"
	case StageDelivered:
		return "Delivered"
	default:
		return ""
	}
}

func nextStages(kind string, current Stage) []Stage {
	if current == "" {
		return []Stage{}
	}
	if kind == KindBespoke {
		switch current {
		case StagePaymentConfirmed:
			return []Stage{StageMeasurementsApproved}
		case StageMeasurementsApproved:
			return []Stage{StageInProduction}
		case StageInProduction:
			return []Stage{StageQualityCheck}
		}
	} else {
		switch current {
		case StagePaymentConfirmed:
			return []Stage{StagePreparingOrder}
		case StagePreparingOrder:
			return []Stage{StageQualityCheck}
		}
	}

	switch current {
	case StageQualityCheck:
		return []Stage{StageReadyForDispatch}
	case StageReadyForDispatch:
		return []Stage{StageDispatched}
	case StageDispatched:
		return []Stage{StageDelivered}
	default:
		return []Stage{}
	}
}

func legacyOrderStatus(stage Stage) string {
	switch stage {
	case StagePaymentConfirmed:
		return "paid"
	case StageDispatched:
		return "shipped"
	case StageDelivered:
		return "delivered"
	default:
		return "in_progress"
	}
}

func parseOptionalDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(dateFormat, value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

// parseOptionalTime reads the handover moment. Parsed in the server's location
// rather than UTC: the admin is typing the local time the parcel changed hands.
func parseOptionalTime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation(dateTimeFormat, value, time.Local)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func validCarrierFields(name string, number string, rawURL string) bool {
	if len(name) > maximumCarrierFieldLength || len(number) > maximumCarrierFieldLength ||
		len(rawURL) > 1000 {
		return false
	}
	if rawURL == "" {
		return true
	}
	parsedURL, err := url.ParseRequestURI(rawURL)
	return err == nil && (parsedURL.Scheme == "https" || parsedURL.Scheme == "http") &&
		parsedURL.Host != ""
}

func containsStage(stages []Stage, candidate Stage) bool {
	for _, stage := range stages {
		if stage == candidate {
			return true
		}
	}
	return false
}

func newTrackingToken() (string, error) {
	tokenBytes := make([]byte, trackingTokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}
