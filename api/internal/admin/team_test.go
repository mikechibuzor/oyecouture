package admin

import (
	"context"
	"errors"
	"testing"
	"time"
)

// teamStoreStub is an in-memory TeamStore. Members are keyed by id so the
// lockout guards can be exercised against realistic combinations.
type teamStoreStub struct {
	members         map[string]TeamMember
	passwordHashes  map[string]string
	deletedSessions []string
	updated         bool
	deleted         bool
}

func newTeamStoreStub() *teamStoreStub {
	return &teamStoreStub{
		members:        map[string]TeamMember{},
		passwordHashes: map[string]string{},
	}
}

func (stub *teamStoreStub) add(id string, role Role, active bool) {
	stub.members[id] = TeamMember{
		ID: id, Email: id + "@example.com", Role: role, Active: active,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func (stub *teamStoreStub) ListMembers(context.Context) ([]TeamMember, error) {
	members := make([]TeamMember, 0, len(stub.members))
	for _, member := range stub.members {
		members = append(members, member)
	}

	return members, nil
}

func (stub *teamStoreStub) MemberByID(_ context.Context, id string) (TeamMember, error) {
	member, ok := stub.members[id]
	if !ok {
		return TeamMember{}, ErrTeamNotFound
	}

	return member, nil
}

func (stub *teamStoreStub) CreateMember(
	_ context.Context, email string, _ string, role Role,
) (TeamMember, error) {
	return TeamMember{ID: "new", Email: email, Role: role, Active: true}, nil
}

func (stub *teamStoreStub) UpdateMember(
	_ context.Context, id string, role Role, active bool,
) (TeamMember, error) {
	stub.updated = true
	member := stub.members[id]
	member.Role = role
	member.Active = active
	stub.members[id] = member

	return member, nil
}

func (stub *teamStoreStub) DeleteMember(_ context.Context, id string) error {
	stub.deleted = true
	delete(stub.members, id)

	return nil
}

func (stub *teamStoreStub) CountActiveOwners(context.Context) (int, error) {
	owners := 0
	for _, member := range stub.members {
		if member.Role == RoleOwner && member.Active {
			owners++
		}
	}

	return owners, nil
}

func (stub *teamStoreStub) PasswordHashByID(_ context.Context, id string) (string, error) {
	hash, ok := stub.passwordHashes[id]
	if !ok {
		return "", ErrTeamNotFound
	}

	return hash, nil
}

func (stub *teamStoreStub) UpdatePassword(_ context.Context, id string, hash string) error {
	stub.passwordHashes[id] = hash

	return nil
}

func (stub *teamStoreStub) DeleteSessionsForUser(_ context.Context, id string) error {
	stub.deletedSessions = append(stub.deletedSessions, id)

	return nil
}

// TestTeamGuardsPreventLockout is the point of this feature's safety rules.
// Every case here, if allowed, ends with nobody able to manage the team except
// through a shell on the server.
func TestTeamGuardsPreventLockout(t *testing.T) {
	t.Parallel()

	t.Run("an owner cannot demote themselves", func(t *testing.T) {
		t.Parallel()

		store := newTeamStoreStub()
		store.add("me", RoleOwner, true)
		store.add("other", RoleOwner, true)

		_, err := NewTeamService(store).Update(context.Background(), "me", "me", RoleStaff, true)
		if !errors.Is(err, ErrSelfTarget) {
			t.Fatalf("error = %v, want ErrSelfTarget", err)
		}
		if store.updated {
			t.Fatal("the account was changed anyway")
		}
	})

	t.Run("an owner cannot delete themselves", func(t *testing.T) {
		t.Parallel()

		store := newTeamStoreStub()
		store.add("me", RoleOwner, true)
		store.add("other", RoleOwner, true)

		err := NewTeamService(store).Delete(context.Background(), "me", "me")
		if !errors.Is(err, ErrSelfTarget) {
			t.Fatalf("error = %v, want ErrSelfTarget", err)
		}
		if store.deleted {
			t.Fatal("the account was deleted anyway")
		}
	})

	t.Run("the last owner cannot be demoted", func(t *testing.T) {
		t.Parallel()

		store := newTeamStoreStub()
		store.add("me", RoleStaff, true)
		store.add("lastOwner", RoleOwner, true)

		_, err := NewTeamService(store).
			Update(context.Background(), "me", "lastOwner", RoleStaff, true)
		if !errors.Is(err, ErrLastOwnerLocked) {
			t.Fatalf("error = %v, want ErrLastOwnerLocked", err)
		}
	})

	t.Run("the last owner cannot be deactivated", func(t *testing.T) {
		t.Parallel()

		store := newTeamStoreStub()
		store.add("me", RoleStaff, true)
		store.add("lastOwner", RoleOwner, true)

		_, err := NewTeamService(store).
			Update(context.Background(), "me", "lastOwner", RoleOwner, false)
		if !errors.Is(err, ErrLastOwnerLocked) {
			t.Fatalf("error = %v, want ErrLastOwnerLocked", err)
		}
	})

	t.Run("the last owner cannot be deleted", func(t *testing.T) {
		t.Parallel()

		store := newTeamStoreStub()
		store.add("me", RoleStaff, true)
		store.add("lastOwner", RoleOwner, true)

		err := NewTeamService(store).Delete(context.Background(), "me", "lastOwner")
		if !errors.Is(err, ErrLastOwnerLocked) {
			t.Fatalf("error = %v, want ErrLastOwnerLocked", err)
		}
		if store.deleted {
			t.Fatal("the last owner was deleted anyway")
		}
	})

	// The guard is about the last owner, not about owners in general.
	t.Run("an owner can be demoted while another remains", func(t *testing.T) {
		t.Parallel()

		store := newTeamStoreStub()
		store.add("me", RoleOwner, true)
		store.add("other", RoleOwner, true)

		if _, err := NewTeamService(store).
			Update(context.Background(), "me", "other", RoleStaff, true); err != nil {
			t.Fatalf("demoting a second owner failed: %v", err)
		}
	})
}

// TestRevokingAccessEndsTheSession keeps revoke meaning what it says. Without
// this the revoked admin keeps working until their session happens to expire.
func TestRevokingAccessEndsTheSession(t *testing.T) {
	t.Parallel()

	store := newTeamStoreStub()
	store.add("me", RoleOwner, true)
	store.add("staff", RoleStaff, true)

	if _, err := NewTeamService(store).
		Update(context.Background(), "me", "staff", RoleStaff, false); err != nil {
		t.Fatalf("revoke failed: %v", err)
	}
	if len(store.deletedSessions) != 1 || store.deletedSessions[0] != "staff" {
		t.Fatalf("sessions deleted = %v, want [staff]", store.deletedSessions)
	}
}

// TestChangePasswordRequiresTheCurrentOne stops a stolen session from silently
// taking the account over and locking the real owner out.
func TestChangePasswordRequiresTheCurrentOne(t *testing.T) {
	t.Parallel()

	currentHash, err := HashPassword("correct-horse-battery")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	store := newTeamStoreStub()
	store.add("me", RoleOwner, true)
	store.passwordHashes["me"] = currentHash
	service := NewTeamService(store)

	if err := service.ChangePassword(
		context.Background(), "me", "not-the-password", "a-brand-new-secret",
	); !errors.Is(err, ErrWrongPassword) {
		t.Fatalf("error = %v, want ErrWrongPassword", err)
	}
	if store.passwordHashes["me"] != currentHash {
		t.Fatal("the password changed despite the wrong current password")
	}

	// Correct current password, but the new one is too short.
	if err := service.ChangePassword(
		context.Background(), "me", "correct-horse-battery", "short",
	); !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("error = %v, want ErrWeakPassword", err)
	}

	if err := service.ChangePassword(
		context.Background(), "me", "correct-horse-battery", "a-brand-new-secret",
	); err != nil {
		t.Fatalf("valid change failed: %v", err)
	}
	if store.passwordHashes["me"] == currentHash {
		t.Fatal("the password did not change")
	}
}

// TestCreateRejectsBadInput keeps weak passwords and unknown roles out.
func TestCreateRejectsBadInput(t *testing.T) {
	t.Parallel()

	service := NewTeamService(newTeamStoreStub())

	if _, err := service.Create(
		context.Background(), "new@example.com", "correct-horse-battery", Role("superuser"),
	); !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("error = %v, want ErrInvalidRole", err)
	}

	if _, err := service.Create(
		context.Background(), "new@example.com", "short", RoleStaff,
	); !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("error = %v, want ErrWeakPassword", err)
	}
}
