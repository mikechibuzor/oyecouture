package payments

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/email"
)

const referenceRandomBytes = 8

type Service struct {
	repository Repository
	provider   Provider
	emailer    PaidOrderEmailSender
	webBaseURL string
}

type PaidOrderEmailSender interface {
	SendPaidOrderEmails(context.Context, email.Order) error
}

// NewService wires payment persistence, provider calls, and optional email notifications.
func NewService(
	repository Repository,
	provider Provider,
	webBaseURL string,
	emailers ...PaidOrderEmailSender,
) *Service {
	var emailer PaidOrderEmailSender
	if len(emailers) > 0 {
		emailer = emailers[0]
	}

	return &Service{
		repository: repository,
		provider:   provider,
		emailer:    emailer,
		webBaseURL: webBaseURL,
	}
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

	reusableAttempt, found, err := service.repository.ReusableInitializedAttempt(
		requestContext, order.ID, service.provider.Name(), order.TotalKobo, CurrencyNGN,
	)
	if err != nil {
		return InitializeResult{}, err
	}
	if found {
		return reusableAttempt, nil
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
		if errors.Is(err, ErrPaymentInProgress) {
			_ = service.repository.MarkAttemptFailed(requestContext, reference, providerPayload)

			reusableAttempt, found, lookupErr := service.repository.ReusableInitializedAttempt(
				requestContext, order.ID, service.provider.Name(), order.TotalKobo, CurrencyNGN,
			)
			if lookupErr != nil {
				return InitializeResult{}, lookupErr
			}
			if found {
				return reusableAttempt, nil
			}
		}
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
	service.sendPaidOrderEmails(requestContext, paidOrder, verification)

	return VerifyResult{
		Order:     paidOrder,
		Reference: verification.Reference,
		Status:    verification.Status,
		Paid:      true,
		Message:   "Payment confirmed.",
		Channel:   verification.Channel,
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

	paidOrder, err := service.repository.ConfirmSuccessfulPayment(requestContext, "", verification)
	if err == nil {
		service.sendPaidOrderEmails(requestContext, paidOrder, verification)
	}
	return err
}

// sendPaidOrderEmails maps a verified paid order into transactional customer emails.
func (service *Service) sendPaidOrderEmails(
	requestContext context.Context,
	order Order,
	verification Verification,
) {
	if service.emailer == nil {
		return
	}

	paidAt := time.Now()
	if verification.PaidAt != nil {
		paidAt = *verification.PaidAt
	}

	emailOrder := email.Order{
		ID:              order.ID,
		OrderNo:         order.OrderNo,
		CustomerName:    order.CustomerName,
		Email:           order.Email,
		Phone:           order.Phone,
		Address:         order.Address,
		DeliveryZone:    order.DeliveryZone,
		SubtotalKobo:    order.SubtotalKobo,
		DeliveryFeeKobo: order.DeliveryFeeKobo,
		TotalKobo:       order.TotalKobo,
		Items:           emailItems(order.Items),
		Payment: email.Payment{
			Provider:  order.PaymentProvider,
			Reference: verification.Reference,
			PaidAt:    paidAt.Format("2 Jan 2006, 15:04 MST"),
		},
	}

	if emailOrder.Payment.Provider == "" {
		emailOrder.Payment.Provider = service.provider.Name()
	}

	if err := service.emailer.SendPaidOrderEmails(requestContext, emailOrder); err != nil {
		// Payment is already confirmed; email failures are recorded by the email service.
		log.Printf("send paid order emails: %v", err)
	}
}

// emailItems converts payment order items into email line items.
func emailItems(items []Item) []email.Item {
	emailItems := make([]email.Item, 0, len(items))
	for _, item := range items {
		emailItems = append(emailItems, email.Item{
			Name:      item.Name,
			Size:      item.Size,
			Quantity:  item.Quantity,
			PriceKobo: item.PriceKobo,
			LineKobo:  item.LineKobo,
		})
	}

	return emailItems
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
