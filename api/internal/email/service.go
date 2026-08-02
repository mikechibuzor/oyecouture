package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"strings"
)

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
