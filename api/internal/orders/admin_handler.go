package orders

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// OrderLister is the admin service contract exposed over HTTP.
type OrderLister interface {
	List(context.Context, ListFilter) (Page, error)
	Metrics(context.Context, MetricsPeriod) (Metrics, error)
	Detail(context.Context, string) (Detail, error)
	UpdateStatus(context.Context, string, StatusUpdate) (Detail, error)
}

type AdminHandler struct {
	orderLister OrderLister
}

func NewAdminHandler(orderLister OrderLister) *AdminHandler {
	return &AdminHandler{orderLister: orderLister}
}

// List serves the admin orders list. It sits behind RequireSession, so reaching
// this handler already means a valid session.
func (handler *AdminHandler) List(responseWriter http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()

	page, err := handler.orderLister.List(request.Context(), ListFilter{
		Status: Status(queryValues.Get("status")),
		Kind:   Kind(queryValues.Get("kind")),
		Search: queryValues.Get("q"),
		From:   ParseFilterDate(queryValues.Get("from")),
		To:     ParseFilterDate(queryValues.Get("to")),
		Limit:  positiveInt(queryValues.Get("limit")),
		Offset: positiveInt(queryValues.Get("offset")),
	})
	if err != nil {
		log.Printf("list admin orders failed: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, ErrorCodeInternal,
			"the orders could not be loaded")
		return
	}

	writeJSON(responseWriter, http.StatusOK, page)
}

// Metrics serves the summary counts shown above the orders list.
func (handler *AdminHandler) Metrics(responseWriter http.ResponseWriter, request *http.Request) {
	metrics, err := handler.orderLister.Metrics(
		request.Context(), MetricsPeriod(request.URL.Query().Get("period")),
	)
	if err != nil {
		log.Printf("load admin order metrics failed: %v", err)
		writeError(responseWriter, http.StatusInternalServerError, ErrorCodeInternal,
			"the summary could not be loaded")
		return
	}

	writeJSON(responseWriter, http.StatusOK, metrics)
}

// Detail serves one admin order detail view.
func (handler *AdminHandler) Detail(responseWriter http.ResponseWriter, request *http.Request) {
	detail, err := handler.orderLister.Detail(request.Context(), request.PathValue("orderNo"))
	if err != nil {
		handler.writeAdminError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, detail)
}

// UpdateStatus applies a service-approved lifecycle action.
func (handler *AdminHandler) UpdateStatus(responseWriter http.ResponseWriter, request *http.Request) {
	var update StatusUpdate

	decoder := json.NewDecoder(http.MaxBytesReader(responseWriter, request.Body, maximumRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&update); err != nil {
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidBody,
			"the status update could not be read")
		return
	}

	detail, err := handler.orderLister.UpdateStatus(
		request.Context(), request.PathValue("orderNo"), update,
	)
	if err != nil {
		handler.writeAdminError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, detail)
}

// writeAdminError maps admin read/action errors to stable HTTP responses.
func (handler *AdminHandler) writeAdminError(responseWriter http.ResponseWriter, adminError error) {
	switch {
	case errors.Is(adminError, ErrOrderNotFound):
		writeError(responseWriter, http.StatusNotFound, ErrorCodeOrderNotFound,
			"the order could not be found")
	case errors.Is(adminError, ErrInvalidStatusUpdate),
		errors.Is(adminError, ErrInvalidStatusTransition),
		errors.Is(adminError, ErrPaymentMustBeServerSide):
		writeError(responseWriter, http.StatusBadRequest, ErrorCodeInvalidStatus,
			"that status change is not available")
	default:
		log.Printf("admin order action failed: %v", adminError)
		writeError(responseWriter, http.StatusInternalServerError, ErrorCodeInternal,
			"the order could not be loaded")
	}
}

// positiveInt reads a paging parameter, treating anything unparseable as unset
// so the service applies its own default.
func positiveInt(value string) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0
	}

	return parsed
}
