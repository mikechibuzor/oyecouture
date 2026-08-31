package tracking

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/uzor-co/oyecouture/api/internal/email"
)

func TestTransitionUsesBespokeSequenceAndQueuesProductionEmail(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{adminState: testAdminState(KindBespoke, StageMeasurementsApproved)}
	service := NewService(repository, email.ProviderSendByte, "https://oyecouture.test")

	state, err := service.Transition(context.Background(), "OYE-1001", "admin-id", TransitionInput{
		Stage:           StageInProduction,
		ExpectedVersion: 2,
		PublicNote:      "Cutting has started.",
	})
	if err != nil {
		t.Fatalf("Transition() error = %v", err)
	}
	if repository.command.NextStage != StageInProduction {
		t.Fatalf("next stage = %q, want %q", repository.command.NextStage, StageInProduction)
	}
	if repository.command.Notification == nil ||
		repository.command.Notification.Type != email.TypeFulfillmentUpdate {
		t.Fatal("production milestone did not queue its email")
	}
	if state.CurrentStage != StageInProduction {
		t.Fatalf("returned stage = %q, want %q", state.CurrentStage, StageInProduction)
	}
}

func TestTransitionRejectsBespokeStageForReadyToWear(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{adminState: testAdminState(KindReadyToWear, StagePaymentConfirmed)}
	service := NewService(repository, email.ProviderSendByte, "https://oyecouture.test")

	_, err := service.Transition(context.Background(), "OYE-1001", "admin-id", TransitionInput{
		Stage:           StageMeasurementsApproved,
		ExpectedVersion: 2,
	})
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("Transition() error = %v, want ErrInvalidTransition", err)
	}
	if repository.transitionCalls != 0 {
		t.Fatal("invalid transition reached the repository")
	}
}

func TestDispatchCreatesTrackingLinkAndEmail(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{adminState: testAdminState(KindReadyToWear, StageReadyForDispatch)}
	service := NewService(repository, email.ProviderSendByte, "https://oyecouture.test")

	state, err := service.Transition(context.Background(), "OYE-1001", "admin-id", TransitionInput{
		Stage:                 StageDispatched,
		ExpectedVersion:       2,
		CarrierName:           "DHL",
		CarrierTrackingNumber: "DHL-123",
		CarrierTrackingURL:    "https://dhl.example/track/DHL-123",
	})
	if err != nil {
		t.Fatalf("Transition() error = %v", err)
	}
	if len(repository.command.TrackingToken) < 40 {
		t.Fatal("dispatch did not generate a strong tracking token")
	}
	if repository.command.Notification == nil ||
		repository.command.Notification.Type != email.TypeOrderDispatched {
		t.Fatal("dispatch did not queue its email")
	}
	if !strings.HasPrefix(state.TrackingURL, "https://oyecouture.test/track/") {
		t.Fatalf("tracking URL = %q", state.TrackingURL)
	}
}

func TestTransitionRejectsStaleVersion(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{adminState: testAdminState(KindReadyToWear, StagePaymentConfirmed)}
	service := NewService(repository, email.ProviderSendByte, "https://oyecouture.test")

	_, err := service.Transition(context.Background(), "OYE-1001", "admin-id", TransitionInput{
		Stage:           StagePreparingOrder,
		ExpectedVersion: 1,
	})
	if !errors.Is(err, ErrVersionConflict) {
		t.Fatalf("Transition() error = %v, want ErrVersionConflict", err)
	}
}

func testAdminState(kind string, stage Stage) AdminState {
	return AdminState{
		OrderNo:       "OYE-1001",
		Kind:          kind,
		CurrentStage:  stage,
		Version:       2,
		orderID:       "order-id",
		customerName:  "Ada",
		customerEmail: "ada@example.com",
	}
}

type repositoryStub struct {
	adminState      AdminState
	command         TransitionCommand
	transitionCalls int
}

func (repository *repositoryStub) Admin(context.Context, string) (AdminState, error) {
	return repository.adminState, nil
}

func (repository *repositoryStub) Public(context.Context, string) (PublicState, error) {
	return PublicState{}, nil
}

func (repository *repositoryStub) Transition(
	_ context.Context,
	command TransitionCommand,
) error {
	repository.transitionCalls++
	repository.command = command
	repository.adminState.CurrentStage = command.NextStage
	repository.adminState.Version++
	repository.adminState.trackingToken = command.TrackingToken
	repository.adminState.Shipment = Shipment{
		CarrierName:           command.CarrierName,
		CarrierTrackingNumber: command.CarrierTrackingNumber,
		CarrierTrackingURL:    command.CarrierTrackingURL,
	}
	return nil
}
