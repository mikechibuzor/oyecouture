package contact

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"net/mail"
	"strings"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/email"
)

const (
	maximumNameLength    = 120
	maximumEmailLength   = 254
	maximumTopicLength   = 80
	maximumMessageLength = 3000
)

var (
	ErrInvalidSubmission = errors.New("invalid contact submission")
	ErrSenderUnavailable = errors.New("contact email sender unavailable")
)

type Service struct {
	provider email.Provider
	from     string
	// Resolved per send rather than held, so changing the address in the admin takes effect
	// on the next enquiry instead of the next deploy.
	recipient func(context.Context) (string, error)
	// Optional. Replaces the display name in the From line; nil keeps the configured header.
	fromOverride func(context.Context, string) string
	now          func() time.Time
}

// WithFromOverride lets the From line carry the name the studio set in the admin. It can only
// replace the display name, never stop an enquiry being sent.
func (service *Service) WithFromOverride(override func(context.Context, string) string) *Service {
	service.fromOverride = override

	return service
}

func NewService(
	provider email.Provider,
	from string,
	recipient func(context.Context) (string, error),
) *Service {
	return &Service{
		provider:  provider,
		from:      strings.TrimSpace(from),
		recipient: recipient,
		now:       time.Now,
	}
}

func (service *Service) Submit(ctx context.Context, submission Submission) error {
	cleanSubmission, err := cleanAndValidate(submission)
	if err != nil {
		return err
	}
	if service.provider == nil || service.from == "" || service.recipient == nil {
		return ErrSenderUnavailable
	}

	// An enquiry with nowhere to go is refused rather than dropped, so the customer is told
	// to try another way instead of believing it was sent.
	notifyEmail, err := service.recipient(ctx)
	if err != nil {
		return fmt.Errorf("resolve contact recipient: %w", err)
	}
	if strings.TrimSpace(notifyEmail) == "" {
		return ErrSenderUnavailable
	}

	fromHeader := service.from
	if service.fromOverride != nil {
		if composed := strings.TrimSpace(service.fromOverride(ctx, service.from)); composed != "" {
			fromHeader = composed
		}
	}

	_, _, err = service.provider.Send(ctx, email.Message{
		From:           fromHeader,
		To:             strings.TrimSpace(notifyEmail),
		Subject:        contactSubject(cleanSubmission),
		HTML:           contactHTML(cleanSubmission),
		Text:           contactText(cleanSubmission),
		IdempotencyKey: contactIdempotencyKey(cleanSubmission, service.now()),
	})
	if err != nil {
		return fmt.Errorf("send contact email: %w", err)
	}

	return nil
}

func cleanAndValidate(submission Submission) (Submission, error) {
	cleanSubmission := Submission{
		Name:    strings.TrimSpace(submission.Name),
		Email:   strings.TrimSpace(submission.Email),
		Topic:   strings.TrimSpace(submission.Topic),
		Rating:  submission.Rating,
		Message: strings.TrimSpace(submission.Message),
	}
	if cleanSubmission.Topic == "" {
		cleanSubmission.Topic = "General"
	}

	if cleanSubmission.Name == "" || len(cleanSubmission.Name) > maximumNameLength ||
		cleanSubmission.Email == "" || len(cleanSubmission.Email) > maximumEmailLength ||
		cleanSubmission.Message == "" || len(cleanSubmission.Message) > maximumMessageLength ||
		len(cleanSubmission.Topic) > maximumTopicLength ||
		strings.ContainsAny(cleanSubmission.Name, "\r\n") ||
		strings.ContainsAny(cleanSubmission.Topic, "\r\n") ||
		cleanSubmission.Rating < 0 || cleanSubmission.Rating > 5 {
		return Submission{}, ErrInvalidSubmission
	}
	if _, err := mail.ParseAddress(cleanSubmission.Email); err != nil {
		return Submission{}, ErrInvalidSubmission
	}

	return cleanSubmission, nil
}

func contactSubject(submission Submission) string {
	return fmt.Sprintf("Oyecouture contact: %s from %s", submission.Topic, submission.Name)
}

func contactHTML(submission Submission) string {
	escapedMessage := html.EscapeString(submission.Message)
	escapedMessage = strings.ReplaceAll(escapedMessage, "\n", "<br>")

	return fmt.Sprintf(`
<h2>New Oyecouture contact message</h2>
<p><strong>Name:</strong> %s</p>
<p><strong>Email:</strong> %s</p>
<p><strong>Topic:</strong> %s</p>
<p><strong>Rating:</strong> %d/5</p>
<hr>
<p>%s</p>
`, html.EscapeString(submission.Name), html.EscapeString(submission.Email), html.EscapeString(submission.Topic), submission.Rating, escapedMessage)
}

func contactText(submission Submission) string {
	return fmt.Sprintf(`New Oyecouture contact message

Name: %s
Email: %s
Topic: %s
Rating: %d/5

%s
`, submission.Name, submission.Email, submission.Topic, submission.Rating, submission.Message)
}

func contactIdempotencyKey(submission Submission, submittedAt time.Time) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf(
		"%s|%s|%s|%d|%s|%d",
		submission.Name,
		submission.Email,
		submission.Topic,
		submission.Rating,
		submission.Message,
		submittedAt.UnixNano(),
	)))

	return "contact_" + hex.EncodeToString(hash[:16])
}
