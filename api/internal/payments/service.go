package payments

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

const referenceRandomBytes = 8

type Service struct {
	repository Repository
	provider   Provider
	webBaseURL string
}

func NewService(repository Repository, provider Provider, webBaseURL string) *Service {
	return &Service{repository: repository, provider: provider, webBaseURL: webBaseURL}
}

func (service *Service) Summary(requestContext context.Context, payToken string) (Summary, error) {
	order, err := service.repository.OrderByPayToken(requestContext, strings.TrimSpace(payToken))
	if err != nil {
		return Summary{}, err
	}

	return Summary{
		Order:   order,
		Payable: order.Status == OrderStatusAwaitingPayment,
		Message: paymentMessage(order.Status),
	}, nil
}

func (service *Service) Initialize(
	requestContext context.Context,
	payToken string,
) (InitializeResult, error) {
	if service.provider == nil {
		return InitializeResult{}, ErrPaymentUnavailable
	}

	order, err := service.repository.OrderByPayToken(requestContext, strings.TrimSpace(payToken))
	if err != nil {
		return InitializeResult{}, err
	}
	if order.Status == OrderStatusPaid {
		return InitializeResult{}, ErrPaymentAlreadyPaid
	}
	if order.Status != OrderStatusAwaitingPayment {
		return InitializeResult{}, ErrOrderNotPayable
	}

	reference, err := newReference(order.OrderNo)
	if err != nil {
		return InitializeResult{}, fmt.Errorf("generate payment reference: %w", err)
	}

	attempt := Attempt{
		OrderID:    order.ID,
		Provider:   service.provider.Name(),
		Reference:  reference,
		AmountKobo: order.TotalKobo,
		Currency:   CurrencyNGN,
		Status:     StatusInitializing,
	}
	if err := service.repository.CreateAttempt(requestContext, attempt); err != nil {
		return InitializeResult{}, err
	}

	result, providerPayload, err := service.provider.Initialize(requestContext, InitializeRequest{
		Email:       order.Email,
		AmountKobo:  order.TotalKobo,
		Currency:    CurrencyNGN,
		Reference:   reference,
		CallbackURL: callbackURL(service.webBaseURL, order.PayToken),
		OrderNo:     order.OrderNo,
		Customer:    order.CustomerName,
	})
	if err != nil {
		_ = service.repository.MarkAttemptFailed(requestContext, reference, providerPayload)
		return InitializeResult{}, err
	}
	if result.Reference != reference {
		_ = service.repository.MarkAttemptFailed(requestContext, reference, providerPayload)
		return InitializeResult{}, ErrReferenceMismatch
	}
	if err := service.repository.MarkAttemptInitialized(
		requestContext, reference, result, providerPayload,
	); err != nil {
		return InitializeResult{}, err
	}

	return result, nil
}

func (service *Service) Verify(
	requestContext context.Context,
	payToken string,
	reference string,
) (VerifyResult, error) {
	if service.provider == nil {
		return VerifyResult{}, ErrPaymentUnavailable
	}

	order, err := service.repository.OrderByPayToken(requestContext, strings.TrimSpace(payToken))
	if err != nil {
		return VerifyResult{}, err
	}

	verification, err := service.provider.Verify(requestContext, strings.TrimSpace(reference))
	if err != nil {
		return VerifyResult{}, err
	}
	if verification.Reference != strings.TrimSpace(reference) {
		return VerifyResult{}, ErrReferenceMismatch
	}

	if verification.Status != StatusSuccess {
		if err := service.repository.RecordVerification(requestContext, order.ID, verification); err != nil {
			return VerifyResult{}, err
		}
		return VerifyResult{
			Order:     order,
			Reference: verification.Reference,
			Status:    verification.Status,
			Message:   "Payment is not complete yet.",
		}, nil
	}

	paidOrder, err := service.repository.ConfirmSuccessfulPayment(
		requestContext, order.ID, verification,
	)
	if err != nil {
		return VerifyResult{}, err
	}

	return VerifyResult{
		Order:     paidOrder,
		Reference: verification.Reference,
		Status:    verification.Status,
		Paid:      true,
		Message:   "Payment confirmed.",
	}, nil
}

func (service *Service) ConfirmWebhookReference(
	requestContext context.Context,
	reference string,
) error {
	if service.provider == nil {
		return ErrPaymentUnavailable
	}

	verification, err := service.provider.Verify(requestContext, strings.TrimSpace(reference))
	if err != nil {
		return err
	}
	if verification.Reference != strings.TrimSpace(reference) {
		return ErrReferenceMismatch
	}
	if verification.Status != StatusSuccess {
		return service.repository.RecordVerification(requestContext, "", verification)
	}

	_, err = service.repository.ConfirmSuccessfulPayment(requestContext, "", verification)
	return err
}

func newReference(orderNo string) (string, error) {
	randomBytes := make([]byte, referenceRandomBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return strings.ToUpper(orderNo) + "-" + hex.EncodeToString(randomBytes), nil
}

func callbackURL(webBaseURL string, payToken string) string {
	return strings.TrimRight(webBaseURL, "/") + "/pay/" + payToken + "/callback"
}

func paymentMessage(orderStatus string) string {
	switch orderStatus {
	case OrderStatusAwaitingPayment:
		return "This order is ready for payment."
	case OrderStatusPaid:
		return "This order has already been paid."
	default:
		return "This order is not ready for payment yet."
	}
}
