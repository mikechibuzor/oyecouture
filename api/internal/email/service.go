package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"
)

const fulfillmentWorkerInterval = 10 * time.Second

type Service struct {
	repository Repository
	provider   Provider
	from       string
	siteURL    string
}

// NewService wires the email repository and provider into the notification service.
func NewService(repository Repository, provider Provider, from string, siteURL string) *Service {
	return &Service{
		repository: repository,
		provider:   provider,
		from:       strings.TrimSpace(from),
		siteURL:    strings.TrimRight(siteURL, "/"),
	}
}

// RunFulfillmentWorker delivers due order-progress emails until shutdown.
func (service *Service) RunFulfillmentWorker(requestContext context.Context) {
	if service == nil || service.provider == nil || service.repository == nil || service.from == "" {
		return
	}

	service.drainFulfillmentQueue(requestContext)
	ticker := time.NewTicker(fulfillmentWorkerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-requestContext.Done():
			return
		case <-ticker.C:
			service.drainFulfillmentQueue(requestContext)
		}
	}
}

// drainFulfillmentQueue processes every currently due email without overlapping claims.
func (service *Service) drainFulfillmentQueue(requestContext context.Context) {
	for requestContext.Err() == nil {
		event, found, err := service.repository.ClaimNext(requestContext)
		if err != nil {
			log.Printf("claim fulfillment email: %v", err)
			return
		}
		if !found {
			return
		}

		service.deliverQueuedEvent(requestContext, event)
	}
}

// deliverQueuedEvent renders and sends one leased event, recording either outcome.
func (service *Service) deliverQueuedEvent(
	requestContext context.Context,
	event QueuedEvent,
) {
	var payload FulfillmentPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		_ = service.repository.MarkFailed(
			requestContext, event.IdempotencyKey, "decode fulfillment payload: "+err.Error(), nil,
		)
		return
	}

	htmlBody, textBody, err := renderFulfillmentEmail(event.Type, payload)
	if err != nil {
		_ = service.repository.MarkFailed(requestContext, event.IdempotencyKey, err.Error(), nil)
		return
	}

	result, providerPayload, err := service.provider.Send(requestContext, Message{
		From:           service.from,
		To:             event.RecipientEmail,
		Subject:        event.Subject,
		HTML:           htmlBody,
		Text:           textBody,
		IdempotencyKey: event.IdempotencyKey,
	})
	if err != nil {
		_ = service.repository.MarkFailed(
			requestContext, event.IdempotencyKey, err.Error(), providerPayload,
		)
		return
	}

	if err := service.repository.MarkSent(
		requestContext, event.IdempotencyKey, result, providerPayload,
	); err != nil {
		log.Printf("mark fulfillment email sent: %v", err)
	}
}

// SendPaidOrderEmails sends the customer confirmation and receipt after verified payment.
func (service *Service) SendPaidOrderEmails(requestContext context.Context, order Order) error {
	if service == nil || service.provider == nil || service.repository == nil || service.from == "" {
		return nil
	}
	if strings.TrimSpace(order.Email) == "" {
		return nil
	}

	if err := service.send(requestContext, order, TypeOrderConfirmation); err != nil {
		log.Printf("send order confirmation email: %v", err)
	}
	if err := service.send(requestContext, order, TypePaymentReceipt); err != nil {
		log.Printf("send payment receipt email: %v", err)
	}

	return nil
}

// send reserves a single email event, renders it, and records the provider result.
func (service *Service) send(
	requestContext context.Context,
	order Order,
	emailType string,
) error {
	subject := emailSubject(emailType, order)
	idempotencyKey := emailIdempotencyKey(emailType, order)
	event := Event{
		OrderID:        order.ID,
		Type:           emailType,
		Provider:       service.provider.Name(),
		IdempotencyKey: idempotencyKey,
		RecipientEmail: order.Email,
		Subject:        subject,
	}

	shouldSend, err := service.repository.Reserve(requestContext, event)
	if err != nil {
		return err
	}
	if !shouldSend {
		return nil
	}

	htmlBody, textBody, err := renderEmail(emailType, service.siteURL, order)
	if err != nil {
		_ = service.repository.MarkFailed(requestContext, idempotencyKey, err.Error(), nil)
		return err
	}

	result, providerPayload, err := service.provider.Send(requestContext, Message{
		From:           service.from,
		To:             order.Email,
		Subject:        subject,
		HTML:           htmlBody,
		Text:           textBody,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		_ = service.repository.MarkFailed(requestContext, idempotencyKey, err.Error(), providerPayload)
		return err
	}

	return service.repository.MarkSent(requestContext, idempotencyKey, result, providerPayload)
}

// emailSubject returns the customer-facing subject for an order email type.
func emailSubject(emailType string, order Order) string {
	switch emailType {
	case TypePaymentReceipt:
		return fmt.Sprintf("Payment received for %s", order.OrderNo)
	default:
		return fmt.Sprintf("Order confirmed - %s", order.OrderNo)
	}
}

// emailIdempotencyKey makes each paid-order email safe to retry.
func emailIdempotencyKey(emailType string, order Order) string {
	switch emailType {
	case TypePaymentReceipt:
		return "order:" + order.OrderNo + ":payment-receipt"
	default:
		return "order:" + order.OrderNo + ":confirmation"
	}
}

// renderEmail renders the HTML and plain-text bodies for an order email type.
func renderEmail(emailType string, siteURL string, order Order) (string, string, error) {
	view := newEmailView(siteURL, order)
	var htmlTemplate string
	switch emailType {
	case TypePaymentReceipt:
		htmlTemplate = paymentReceiptHTML
	default:
		htmlTemplate = orderConfirmationHTML
	}

	parsedTemplate, err := template.New(emailType).Parse(htmlTemplate)
	if err != nil {
		return "", "", err
	}

	var htmlBuffer bytes.Buffer
	if err := parsedTemplate.Execute(&htmlBuffer, view); err != nil {
		return "", "", err
	}

	return htmlBuffer.String(), renderText(emailType, view), nil
}
