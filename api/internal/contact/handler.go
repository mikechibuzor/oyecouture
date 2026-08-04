package contact

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/uzor-co/oyecouture/api/internal/ratelimit"
)

type ErrorCode string

const (
	ErrorCodeInvalidBody       ErrorCode = "invalid_body"
	ErrorCodeInvalidSubmission ErrorCode = "invalid_submission"
	ErrorCodeEmailUnavailable  ErrorCode = "email_unavailable"
	ErrorCodeInternal          ErrorCode = "internal_error"
	ErrorCodeTooManyRequests   ErrorCode = "too_many_requests"

	maximumRequestBytes = 8 << 10

	// The form is public and every accepted submission sends an email, so the
	// abuse here is volume rather than guessing. A real enquiry is one message;
	// this leaves room to correct a typo and resend without ever letting a
	// script flood the owner's inbox, burn the SendByte quota, or damage the
	// sending domain's reputation — which would take order emails down with it.
	submissionWindow   = time.Hour
	submissionsPerHour = 5
)

type Submitter interface {
	Submit(context.Context, Submission) error
}

type Handler struct {
	submitter Submitter
	limiter   *ratelimit.Limiter
}

func NewHandler(submitter Submitter) *Handler {
	return &Handler{
		submitter: submitter,
		limiter:   ratelimit.New(submissionWindow, submissionsPerHour),
	}
}

type errorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

func (handler *Handler) Submit(responseWriter http.ResponseWriter, request *http.Request) {
	address := ratelimit.ClientAddress(request)
	if !handler.limiter.Allow(address) {
		writeError(responseWriter, http.StatusTooManyRequests, ErrorCodeTooManyRequests,
			"too many messages, please try again later")
		return
	}

	var submission Submission

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, maximumRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&submission); err != nil {
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidBody,
			"the message could not be read")
		return
	}

	if err := handler.submitter.Submit(request.Context(), submission); err != nil {
		handler.writeServiceError(responseWriter, err)
		return
	}

	handler.limiter.Record(address)

	writeJSON(responseWriter, http.StatusCreated, Response{Received: true})
}

func (handler *Handler) writeServiceError(responseWriter http.ResponseWriter, serviceError error) {
	switch {
	case errors.Is(serviceError, ErrInvalidSubmission):
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidSubmission,
			"please add a valid name, email and message")
	case errors.Is(serviceError, ErrSenderUnavailable):
		writeError(responseWriter, http.StatusServiceUnavailable, ErrorCodeEmailUnavailable,
			"messages are temporarily unavailable")
	default:
		log.Printf("submit contact message failed: %v", serviceError)
		writeError(responseWriter, http.StatusInternalServerError, ErrorCodeInternal,
			"the message could not be sent")
	}
}

func writeError(
	responseWriter http.ResponseWriter,
	statusCode int,
	errorCode ErrorCode,
	errorMessage string,
) {
	writeJSON(responseWriter, statusCode, errorResponse{
		Error: errorDetail{Code: errorCode, Message: errorMessage},
	})
}

func writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	if err := json.NewEncoder(responseWriter).Encode(payload); err != nil {
		log.Printf("encode contact response: %v", err)
	}
}
