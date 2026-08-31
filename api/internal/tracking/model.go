package tracking

import "time"

// Stage is the customer-facing fulfillment lifecycle. Payment remains owned by
// the payment service; this timeline starts once payment has been verified.
type Stage string

const (
	StagePaymentConfirmed     Stage = "payment_confirmed"
	StageMeasurementsApproved Stage = "measurements_approved"
	StagePreparingOrder       Stage = "preparing_order"
	StageInProduction         Stage = "in_production"
	StageQualityCheck         Stage = "quality_check"
	StageReadyForDispatch     Stage = "ready_for_dispatch"
	StageDispatched           Stage = "dispatched"
	StageDelivered            Stage = "delivered"
)

const (
	KindBespoke     = "bespoke"
	KindReadyToWear = "ready-to-wear"
)

// Event is one immutable entry in the order timeline.
type Event struct {
	ID         string    `json:"id"`
	Stage      Stage     `json:"stage"`
	Label      string    `json:"label"`
	PublicNote string    `json:"publicNote,omitempty"`
	ActorEmail string    `json:"actorEmail,omitempty"`
	OccurredAt time.Time `json:"occurredAt"`
}

// Item is the deliberately small line-item shape exposed by public tracking.
type Item struct {
	Name     string `json:"name"`
	Size     string `json:"size,omitempty"`
	Quantity int    `json:"quantity"`
}

// Shipment contains optional courier metadata entered at dispatch.
type Shipment struct {
	CarrierName           string `json:"carrierName,omitempty"`
	CarrierTrackingNumber string `json:"carrierTrackingNumber,omitempty"`
	CarrierTrackingURL    string `json:"carrierTrackingUrl,omitempty"`
}

// AdminState is the order-progress model used by the authenticated dashboard.
type AdminState struct {
	OrderNo               string     `json:"orderNo"`
	Kind                  string     `json:"kind"`
	CurrentStage          Stage      `json:"currentStage,omitempty"`
	CurrentStageLabel     string     `json:"currentStageLabel,omitempty"`
	Version               int        `json:"version"`
	EstimatedDeliveryDate *time.Time `json:"estimatedDeliveryDate,omitempty"`
	Shipment              Shipment   `json:"shipment"`
	TrackingURL           string     `json:"trackingUrl,omitempty"`
	Events                []Event    `json:"events"`
	AllowedNextStages     []Stage    `json:"allowedNextStages"`

	orderID       string
	customerName  string
	customerEmail string
	trackingToken string
}

// PublicState omits contact, address, measurements, payment, and admin identity.
type PublicState struct {
	OrderNo               string     `json:"orderNo"`
	Kind                  string     `json:"kind"`
	CurrentStage          Stage      `json:"currentStage"`
	CurrentStageLabel     string     `json:"currentStageLabel"`
	EstimatedDeliveryDate *time.Time `json:"estimatedDeliveryDate,omitempty"`
	Shipment              Shipment   `json:"shipment"`
	Items                 []Item     `json:"items"`
	Events                []Event    `json:"events"`
}

// TransitionInput is the untrusted admin request body.
type TransitionInput struct {
	Stage                 Stage  `json:"stage"`
	ExpectedVersion       int    `json:"expectedVersion"`
	PublicNote            string `json:"publicNote"`
	EstimatedDeliveryDate string `json:"estimatedDeliveryDate"`
	CarrierName           string `json:"carrierName"`
	CarrierTrackingNumber string `json:"carrierTrackingNumber"`
	CarrierTrackingURL    string `json:"carrierTrackingUrl"`
}
