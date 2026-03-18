package workspace

import (
	"strings"
	"time"

	"asana-clone-backend/internal/domain/shared"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleMember:
		return true
	}
	return false
}

type Workspace struct {
	ID          shared.ID
	Name        string
	Description string
	OwnerID     shared.ID
	Members     []WorkspaceMember
	CreatedAt   time.Time
}

type WorkspaceMember struct {
	UserID      shared.ID
	WorkspaceID shared.ID
	Role        Role
	JoinedAt    time.Time
}

func NewWorkspace(name, description string, ownerID shared.ID) (*Workspace, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, shared.NewValidationError("INVALID_NAME", "workspace name must not be empty", "name")
	}

	now := time.Now()
	ws := &Workspace{
		ID:          shared.NewID(),
		Name:        name,
		Description: strings.TrimSpace(description),
		OwnerID:     ownerID,
		CreatedAt:   now,
	}

	ws.Members = []WorkspaceMember{
		{
			UserID:      ownerID,
			WorkspaceID: ws.ID,
			Role:        RoleOwner,
			JoinedAt:    now,
		},
	}

	return ws, nil
}

func (w *Workspace) AddMember(userID shared.ID, role Role) error {
	if !role.IsValid() {
		return shared.NewValidationError("INVALID_ROLE", "invalid workspace role", "role")
	}
	if role == RoleOwner {
		return shared.NewDomainError("INVALID_ROLE", "cannot assign owner role via AddMember")
	}
	if w.IsMember(userID) {
		return shared.NewDomainError("ALREADY_MEMBER", "user is already a member of this workspace")
	}

	w.Members = append(w.Members, WorkspaceMember{
		UserID:      userID,
		WorkspaceID: w.ID,
		Role:        role,
		JoinedAt:    time.Now(),
	})
	return nil
}

func (w *Workspace) RemoveMember(userID shared.ID) error {
	if userID == w.OwnerID {
		return shared.NewDomainError("CANNOT_REMOVE_OWNER", "cannot remove the workspace owner")
	}

	for i, m := range w.Members {
		if m.UserID == userID {
			w.Members = append(w.Members[:i], w.Members[i+1:]...)
			return nil
		}
	}
	return shared.NewDomainError("NOT_MEMBER", "user is not a member of this workspace")
}

func (w *Workspace) IsMember(userID shared.ID) bool {
	for _, m := range w.Members {
		if m.UserID == userID {
			return true
		}
	}
	return false
}
