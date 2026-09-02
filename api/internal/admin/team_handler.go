package admin

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type TeamHandler struct {
	service *TeamService
}

func NewTeamHandler(service *TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

// RequireOwner gates the routes only a super admin may reach. It runs inside
// RequireSession, so a session is already present; a missing one here means the
// route was wired without authentication and must fail closed rather than
// treating an anonymous caller as an owner.
func RequireOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		session, ok := SessionFrom(request.Context())
		if !ok || session.User.Role != RoleOwner {
			writeJSON(responseWriter, http.StatusForbidden,
				map[string]string{"error": "only an owner can manage the team"})
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}

type teamMemberRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Active   *bool  `json:"active"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (handler *TeamHandler) List(responseWriter http.ResponseWriter, request *http.Request) {
	members, err := handler.service.List(request.Context())
	if err != nil {
		handler.writeError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, map[string]any{"members": members})
}

func (handler *TeamHandler) Create(responseWriter http.ResponseWriter, request *http.Request) {
	var payload teamMemberRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest,
			map[string]string{"error": "that request could not be read"})
		return
	}

	member, err := handler.service.Create(
		request.Context(), payload.Email, payload.Password, payload.Role,
	)
	if err != nil {
		handler.writeError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusCreated, map[string]any{"member": member})
}

func (handler *TeamHandler) Update(responseWriter http.ResponseWriter, request *http.Request) {
	session, _ := SessionFrom(request.Context())

	var payload teamMemberRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest,
			map[string]string{"error": "that request could not be read"})
		return
	}

	// Absent means unchanged rather than false, so a role edit cannot revoke
	// someone's access as a side effect of omitting a field.
	active := true
	if payload.Active != nil {
		active = *payload.Active
	}

	member, err := handler.service.Update(
		request.Context(), session.User.ID, request.PathValue("id"), payload.Role, active,
	)
	if err != nil {
		handler.writeError(responseWriter, err)
		return
	}

	writeJSON(responseWriter, http.StatusOK, map[string]any{"member": member})
}

func (handler *TeamHandler) Delete(responseWriter http.ResponseWriter, request *http.Request) {
	session, _ := SessionFrom(request.Context())

	if err := handler.service.Delete(
		request.Context(), session.User.ID, request.PathValue("id"),
	); err != nil {
		handler.writeError(responseWriter, err)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}

// ChangePassword is the one team route any signed-in admin may call: it only
// ever acts on the caller's own account.
func (handler *TeamHandler) ChangePassword(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	session, ok := SessionFrom(request.Context())
	if !ok {
		writeJSON(responseWriter, http.StatusUnauthorized,
			map[string]string{"error": "not signed in"})
		return
	}

	var payload changePasswordRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest,
			map[string]string{"error": "that request could not be read"})
		return
	}

	if err := handler.service.ChangePassword(
		request.Context(), session.User.ID, payload.CurrentPassword, payload.NewPassword,
	); err != nil {
		handler.writeError(responseWriter, err)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}

// writeError maps the service's errors onto status codes. Everything a person
// can fix says what to do; anything else is logged and reported as a server
// error rather than leaking internals to the browser.
func (handler *TeamHandler) writeError(responseWriter http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTeamNotFound):
		writeJSON(responseWriter, http.StatusNotFound,
			map[string]string{"error": ErrTeamNotFound.Error()})
	case errors.Is(err, ErrEmailTaken):
		writeJSON(responseWriter, http.StatusConflict,
			map[string]string{"error": ErrEmailTaken.Error()})
	case errors.Is(err, ErrInvalidRole),
		errors.Is(err, ErrWeakPassword):
		writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrWrongPassword):
		writeJSON(responseWriter, http.StatusUnauthorized,
			map[string]string{"error": ErrWrongPassword.Error()})
	case errors.Is(err, ErrSelfTarget), errors.Is(err, ErrLastOwnerLocked):
		writeJSON(responseWriter, http.StatusConflict, map[string]string{"error": err.Error()})
	default:
		log.Printf("admin team: %v", err)
		writeJSON(responseWriter, http.StatusInternalServerError,
			map[string]string{"error": "that could not be completed"})
	}
}
