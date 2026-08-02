package orders

import (
	"context"
	"log"
	"net/http"
	"strconv"
)

// OrderLister is the admin service contract exposed over HTTP.
type OrderLister interface {
	List(context.Context, ListFilter) (Page, error)
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

// positiveInt reads a paging parameter, treating anything unparseable as unset
// so the service applies its own default.
func positiveInt(value string) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0
	}

	return parsed
}
