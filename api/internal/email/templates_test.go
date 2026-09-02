package email

import (
	"strings"
	"testing"
)

func fulfillmentFixture() FulfillmentPayload {
	return FulfillmentPayload{
		CustomerName:          "Tunde Adeyemi",
		OrderNo:               "OYE-1042",
		StageLabel:            "Dispatched",
		CarrierName:           "DHL Express",
		CarrierTrackingNumber: "7248193055",
		EstimatedDeliveryDate: "5 Sep 2026",
		TrackingURL:           "https://oyecouturebespoke.com/track/token123",
		Items: []Item{
			{Name: "The Sovereign Suit", Size: "L", Quantity: 1, LineKobo: 62000000},
		},
		SubtotalKobo:    62000000,
		DeliveryFeeKobo: 1500000,
		TotalKobo:       63500000,
	}
}

// TestFulfillmentEmailsCarryTheirDesignBlocks pins the parts of the client's
// designs that come from data rather than markup, so a payload field quietly
// dropping out shows up here instead of in a customer's inbox.
func TestFulfillmentEmailsCarryTheirDesignBlocks(t *testing.T) {
	t.Parallel()

	t.Run("dispatched shows the tracking panel", func(t *testing.T) {
		t.Parallel()

		html, _, err := renderFulfillmentEmail(TypeOrderDispatched, fulfillmentFixture())
		if err != nil {
			t.Fatalf("renderFulfillmentEmail: %v", err)
		}

		for _, expected := range []string{
			"DHL Express", "7248193055", "5 Sep 2026",
			"The Sovereign Suit", "Tunde,", "Track your parcel",
		} {
			if !strings.Contains(html, expected) {
				t.Fatalf("dispatched email is missing %q", expected)
			}
		}
	})

	t.Run("delivered shows the handover record", func(t *testing.T) {
		t.Parallel()

		payload := fulfillmentFixture()
		payload.StageLabel = "Delivered"
		payload.SignedForBy = "T. Adeyemi"
		payload.SignedAt = "2 September 2026, 13:42"
		payload.DeliveryLocation = "18 Glover Road, Ikoyi"

		html, _, err := renderFulfillmentEmail(TypeOrderDelivered, payload)
		if err != nil {
			t.Fatalf("renderFulfillmentEmail: %v", err)
		}

		for _, expected := range []string{
			"Signed for by", "T. Adeyemi", "2 September 2026, 13:42", "18 Glover Road, Ikoyi",
		} {
			if !strings.Contains(html, expected) {
				t.Fatalf("delivered email is missing %q", expected)
			}
		}
		// The design's review panel is deliberately absent until a review system
		// exists; a button leading nowhere is worse than no button.
		if strings.Contains(html, "Leave a review") {
			t.Fatal("delivered email offers a review with nowhere to go")
		}
	})

	t.Run("the handover panel is omitted when nothing was recorded", func(t *testing.T) {
		t.Parallel()

		html, _, err := renderFulfillmentEmail(TypeOrderDelivered, fulfillmentFixture())
		if err != nil {
			t.Fatalf("renderFulfillmentEmail: %v", err)
		}
		if strings.Contains(html, "Signed for by") {
			t.Fatal("delivered email shows an empty handover panel")
		}
	})

	// Mail clients block remote images by default, so a hosted logo leaves most
	// recipients with no brand at all. The wordmark is set as type.
	t.Run("no email relies on a hosted image", func(t *testing.T) {
		t.Parallel()

		for _, emailType := range []string{
			TypeFulfillmentUpdate, TypeOrderDispatched, TypeOrderDelivered,
		} {
			html, _, err := renderFulfillmentEmail(emailType, fulfillmentFixture())
			if err != nil {
				t.Fatalf("renderFulfillmentEmail(%s): %v", emailType, err)
			}
			if strings.Contains(html, "<img") {
				t.Fatalf("%s relies on an externally hosted image", emailType)
			}
		}
	})
}

// TestProgressReachesTheRightStep keeps the four-step bar honest: a shipped
// email claiming delivery, or a delivered one still showing it in transit,
// would contradict the message beside it.
func TestProgressReachesTheRightStep(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		emailType string
		wantDone  int
	}{
		{emailType: TypeFulfillmentUpdate, wantDone: 2},
		{emailType: TypeOrderDispatched, wantDone: 3},
		{emailType: TypeOrderDelivered, wantDone: 4},
	}

	for _, testCase := range testCases {
		t.Run(testCase.emailType, func(t *testing.T) {
			t.Parallel()

			steps := progressFor(testCase.emailType)
			if len(steps) != 4 {
				t.Fatalf("progress steps = %d, want 4", len(steps))
			}

			done := 0
			for _, step := range steps {
				if step.Done {
					done++
				}
			}
			if done != testCase.wantDone {
				t.Fatalf("completed steps = %d, want %d", done, testCase.wantDone)
			}
		})
	}
}
