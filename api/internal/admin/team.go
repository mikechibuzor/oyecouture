package admin

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Role is a defined type rather than a bare string, matching tracking.Stage,
// so an arbitrary string cannot be passed where a role is expected.
type Role string

const (
	// RoleOwner is the super admin: manages the team, settings, delivery and
	// prices. RoleStaff runs orders and fulfilment and nothing else.
	RoleOwner Role = "owner"
	RoleStaff Role = "staff"

	minimumPasswordLength = 12
)

var (
	ErrForbidden       = errors.New("not permitted")
	ErrTeamNotFound    = errors.New("admin not found")
	ErrEmailTaken      = errors.New("that email already has an account")
	ErrInvalidRole     = errors.New("role must be owner or staff")
	ErrWeakPassword    = errors.New("password must be at least 12 characters")
	ErrWrongPassword   = errors.New("current password is incorrect")
	ErrSelfTarget      = errors.New("you cannot change your own access")
	ErrLastOwnerLocked = errors.New("the last owner cannot be removed or demoted")
)

// TeamMember is one admin account as the team desk sees it. It deliberately
// carries no password hash: nothing outside authentication should handle one.
type TeamMember struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TeamStore is the persistence contract for managing admin accounts.
type TeamStore interface {
	ListMembers(context.Context) ([]TeamMember, error)
	MemberByID(context.Context, string) (TeamMember, error)
	CreateMember(context.Context, string, string, Role) (TeamMember, error)
	UpdateMember(context.Context, string, Role, bool) (TeamMember, error)
	DeleteMember(context.Context, string) error
	CountActiveOwners(context.Context) (int, error)
	PasswordHashByID(context.Context, string) (string, error)
	UpdatePassword(context.Context, string, string) error
	DeleteSessionsForUser(context.Context, string) error
}

type TeamService struct {
	store TeamStore
}

func NewTeamService(store TeamStore) *TeamService {
	return &TeamService{store: store}
}

func validRole(role Role) bool {
	return role == RoleOwner || role == RoleStaff
}

// List returns every admin account, active and revoked alike.
func (service *TeamService) List(requestContext context.Context) ([]TeamMember, error) {
	return service.store.ListMembers(requestContext)
}

// Create adds an admin. The caller chooses the password, so the new admin
// should change it on first sign-in; that is what ChangePassword exists for.
func (service *TeamService) Create(
	requestContext context.Context,
	email string,
	password string,
	role Role,
) (TeamMember, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return TeamMember{}, ErrTeamNotFound
	}
	if !validRole(role) {
		return TeamMember{}, ErrInvalidRole
	}
	if len(password) < minimumPasswordLength {
		return TeamMember{}, ErrWeakPassword
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return TeamMember{}, err
	}

	return service.store.CreateMember(requestContext, email, passwordHash, role)
}

// Update changes a member's role or active flag.
//
// The guards here are the whole point of the function. An owner who demotes or
// deactivates themselves, or removes the last remaining owner, locks everyone
// out of team management with no way back except a shell on the server.
func (service *TeamService) Update(
	requestContext context.Context,
	actorID string,
	memberID string,
	role Role,
	active bool,
) (TeamMember, error) {
	if !validRole(role) {
		return TeamMember{}, ErrInvalidRole
	}
	if memberID == actorID {
		return TeamMember{}, ErrSelfTarget
	}

	member, err := service.store.MemberByID(requestContext, memberID)
	if err != nil {
		return TeamMember{}, err
	}

	// Losing the last owner is only possible by demoting or deactivating one.
	losingAnOwner := member.Role == RoleOwner && member.Active &&
		(role != RoleOwner || !active)
	if losingAnOwner {
		owners, err := service.store.CountActiveOwners(requestContext)
		if err != nil {
			return TeamMember{}, err
		}
		if owners <= 1 {
			return TeamMember{}, ErrLastOwnerLocked
		}
	}

	updated, err := service.store.UpdateMember(requestContext, memberID, role, active)
	if err != nil {
		return TeamMember{}, err
	}

	// Revoking access has to end their current session too. Without this they
	// keep working until it happens to expire, which is not what revoke means.
	if !active {
		if err := service.store.DeleteSessionsForUser(requestContext, memberID); err != nil {
			return TeamMember{}, err
		}
	}

	return updated, nil
}

// Delete removes an account outright. Prefer deactivating: a deleted admin
// takes their name off the fulfilment events that reference them.
func (service *TeamService) Delete(
	requestContext context.Context,
	actorID string,
	memberID string,
) error {
	if memberID == actorID {
		return ErrSelfTarget
	}

	member, err := service.store.MemberByID(requestContext, memberID)
	if err != nil {
		return err
	}

	if member.Role == RoleOwner && member.Active {
		owners, err := service.store.CountActiveOwners(requestContext)
		if err != nil {
			return err
		}
		if owners <= 1 {
			return ErrLastOwnerLocked
		}
	}

	return service.store.DeleteMember(requestContext, memberID)
}

// ChangePassword lets an admin set their own password.
//
// This exists so the person who created an account is not left holding its
// password forever. Without it, every account the developer seeds keeps a
// secret the developer knows, and an audit trail saying what someone did means
// very little when two people could have been signed in as them.
func (service *TeamService) ChangePassword(
	requestContext context.Context,
	actorID string,
	currentPassword string,
	newPassword string,
) error {
	if len(newPassword) < minimumPasswordLength {
		return ErrWeakPassword
	}

	currentHash, err := service.store.PasswordHashByID(requestContext, actorID)
	if err != nil {
		return err
	}
	// The current password is required so a stolen session cannot quietly
	// change the password and lock the real owner out.
	if !passwordMatches(currentHash, currentPassword) {
		return ErrWrongPassword
	}

	newHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	return service.store.UpdatePassword(requestContext, actorID, newHash)
}
