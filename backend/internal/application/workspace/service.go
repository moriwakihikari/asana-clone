package workspace

import (
	"context"
	"strings"
	"time"

	"asana-clone-backend/internal/domain/shared"
	userDomain "asana-clone-backend/internal/domain/user"
	wsDomain "asana-clone-backend/internal/domain/workspace"
)

// WorkspaceResponse is the public representation of a workspace (list view).
type WorkspaceResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	MemberCount int    `json:"member_count"`
	CreatedAt   string `json:"created_at"`
}

// MemberResponse represents a workspace member with user details.
type MemberResponse struct {
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joined_at"`
}

// WorkspaceDetailResponse is the detailed representation including members.
type WorkspaceDetailResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	OwnerID     string            `json:"owner_id"`
	Members     []*MemberResponse `json:"members"`
	CreatedAt   string            `json:"created_at"`
}

// WorkspaceService handles workspace operations.
type WorkspaceService struct {
	workspaceRepo wsDomain.WorkspaceRepository
	userRepo      userDomain.UserRepository
}

// NewWorkspaceService creates a new WorkspaceService.
func NewWorkspaceService(workspaceRepo wsDomain.WorkspaceRepository, userRepo userDomain.UserRepository) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
	}
}

// Create creates a new workspace owned by ownerID.
func (s *WorkspaceService) Create(ctx context.Context, ownerID shared.ID, name, description string) (*WorkspaceResponse, error) {
	ws, err := wsDomain.NewWorkspace(name, description, ownerID)
	if err != nil {
		return nil, err
	}

	if err := s.workspaceRepo.Save(ctx, ws); err != nil {
		return nil, err
	}

	return toWorkspaceResponse(ws), nil
}

// GetByID retrieves a workspace by ID. The caller must be a member.
func (s *WorkspaceService) GetByID(ctx context.Context, id, callerID shared.ID) (*WorkspaceDetailResponse, error) {
	ws, err := s.workspaceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ws == nil {
		return nil, shared.ErrNotFound
	}

	if !ws.IsMember(callerID) {
		return nil, shared.ErrForbidden
	}

	// Resolve member details.
	members := make([]*MemberResponse, 0, len(ws.Members))
	for _, m := range ws.Members {
		u, err := s.userRepo.FindByID(ctx, m.UserID)
		if err != nil {
			return nil, err
		}
		if u == nil {
			continue
		}
		members = append(members, &MemberResponse{
			UserID:    u.ID.String(),
			Name:      u.Name,
			Email:     u.Email,
			AvatarURL: u.AvatarURL,
			Role:      string(m.Role),
			JoinedAt:  m.JoinedAt.Format(time.RFC3339),
		})
	}

	return &WorkspaceDetailResponse{
		ID:          ws.ID.String(),
		Name:        ws.Name,
		Description: ws.Description,
		OwnerID:     ws.OwnerID.String(),
		Members:     members,
		CreatedAt:   ws.CreatedAt.Format(time.RFC3339),
	}, nil
}

// ListByUser returns all workspaces the user is a member of.
func (s *WorkspaceService) ListByUser(ctx context.Context, userID shared.ID) ([]*WorkspaceResponse, error) {
	workspaces, err := s.workspaceRepo.FindByMemberID(ctx, userID)
	if err != nil {
		return nil, err
	}

	results := make([]*WorkspaceResponse, 0, len(workspaces))
	for _, ws := range workspaces {
		results = append(results, toWorkspaceResponse(ws))
	}
	return results, nil
}

// Update modifies the workspace name and description. Only owner or admin can update.
func (s *WorkspaceService) Update(ctx context.Context, id, callerID shared.ID, name, description string) error {
	ws, err := s.workspaceRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if ws == nil {
		return shared.ErrNotFound
	}

	if !s.isAdminOrOwner(ws, callerID) {
		return shared.ErrForbidden
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return shared.NewValidationError("INVALID_NAME", "workspace name must not be empty", "name")
	}

	ws.Name = name
	ws.Description = strings.TrimSpace(description)

	return s.workspaceRepo.Save(ctx, ws)
}

// Delete removes a workspace. Only the owner can delete.
func (s *WorkspaceService) Delete(ctx context.Context, id, callerID shared.ID) error {
	ws, err := s.workspaceRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if ws == nil {
		return shared.ErrNotFound
	}

	if ws.OwnerID != callerID {
		return shared.ErrForbidden
	}

	return s.workspaceRepo.Delete(ctx, id)
}

// AddMember adds a user to the workspace by email. Only owner or admin can add members.
func (s *WorkspaceService) AddMember(ctx context.Context, workspaceID shared.ID, email string, role string, callerID shared.ID) error {
	ws, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return err
	}
	if ws == nil {
		return shared.ErrNotFound
	}

	if !s.isAdminOrOwner(ws, callerID) {
		return shared.ErrForbidden
	}

	wsRole := wsDomain.Role(role)
	if !wsRole.IsValid() {
		return shared.NewValidationError("INVALID_ROLE", "invalid workspace role", "role")
	}

	// Find user by email.
	u, err := s.userRepo.FindByEmail(ctx, strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		return err
	}
	if u == nil {
		return shared.NewDomainError("USER_NOT_FOUND", "no user found with the given email")
	}

	if err := ws.AddMember(u.ID, wsRole); err != nil {
		return err
	}

	return s.workspaceRepo.Save(ctx, ws)
}

// RemoveMember removes a user from the workspace. Only owner or admin can remove members.
func (s *WorkspaceService) RemoveMember(ctx context.Context, workspaceID, userID, callerID shared.ID) error {
	ws, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return err
	}
	if ws == nil {
		return shared.ErrNotFound
	}

	if !s.isAdminOrOwner(ws, callerID) {
		return shared.ErrForbidden
	}

	if err := ws.RemoveMember(userID); err != nil {
		return err
	}

	return s.workspaceRepo.Save(ctx, ws)
}

// isAdminOrOwner checks whether the caller is the owner or an admin of the workspace.
func (s *WorkspaceService) isAdminOrOwner(ws *wsDomain.Workspace, callerID shared.ID) bool {
	for _, m := range ws.Members {
		if m.UserID == callerID {
			return m.Role == wsDomain.RoleOwner || m.Role == wsDomain.RoleAdmin
		}
	}
	return false
}

func toWorkspaceResponse(ws *wsDomain.Workspace) *WorkspaceResponse {
	return &WorkspaceResponse{
		ID:          ws.ID.String(),
		Name:        ws.Name,
		Description: ws.Description,
		OwnerID:     ws.OwnerID.String(),
		MemberCount: len(ws.Members),
		CreatedAt:   ws.CreatedAt.Format(time.RFC3339),
	}
}
