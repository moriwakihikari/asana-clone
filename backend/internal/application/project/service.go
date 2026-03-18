package project

import (
	"context"
	"time"

	"asana-clone-backend/internal/domain/shared"
	projectDomain "asana-clone-backend/internal/domain/project"
	sectionDomain "asana-clone-backend/internal/domain/section"
	wsDomain "asana-clone-backend/internal/domain/workspace"
)

// ProjectResponse is the public representation of a project.
type ProjectResponse struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	ViewType    string `json:"view_type"`
	IsArchived  bool   `json:"is_archived"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ProjectService handles project operations.
type ProjectService struct {
	projectRepo   projectDomain.ProjectRepository
	sectionRepo   sectionDomain.SectionRepository
	workspaceRepo wsDomain.WorkspaceRepository
}

// NewProjectService creates a new ProjectService.
func NewProjectService(
	projectRepo projectDomain.ProjectRepository,
	sectionRepo sectionDomain.SectionRepository,
	workspaceRepo wsDomain.WorkspaceRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo:   projectRepo,
		sectionRepo:   sectionRepo,
		workspaceRepo: workspaceRepo,
	}
}

// Create creates a new project and auto-creates three default sections (To Do, In Progress, Done).
func (s *ProjectService) Create(
	ctx context.Context,
	workspaceID, callerID shared.ID,
	name, description, color string,
	viewType string,
) (*ProjectResponse, error) {
	// Verify the caller is a member of the workspace.
	ws, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if ws == nil {
		return nil, shared.ErrNotFound
	}
	if !ws.IsMember(callerID) {
		return nil, shared.ErrForbidden
	}

	p, err := projectDomain.NewProject(workspaceID, name, description, color, projectDomain.ViewType(viewType))
	if err != nil {
		return nil, err
	}

	if err := s.projectRepo.Save(ctx, p); err != nil {
		return nil, err
	}

	// Create default sections.
	defaultSections := []string{"To Do", "In Progress", "Done"}
	for i, sName := range defaultSections {
		sec, err := sectionDomain.NewSection(p.ID, sName, i)
		if err != nil {
			return nil, err
		}
		if err := s.sectionRepo.Save(ctx, sec); err != nil {
			return nil, err
		}
	}

	return toProjectResponse(p), nil
}

// GetByID retrieves a project by ID.
func (s *ProjectService) GetByID(ctx context.Context, id shared.ID) (*ProjectResponse, error) {
	p, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, shared.ErrNotFound
	}
	return toProjectResponse(p), nil
}

// ListByWorkspace returns all projects in a workspace.
func (s *ProjectService) ListByWorkspace(ctx context.Context, workspaceID, callerID shared.ID) ([]*ProjectResponse, error) {
	// Verify the caller is a member of the workspace.
	ws, err := s.workspaceRepo.FindByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if ws == nil {
		return nil, shared.ErrNotFound
	}
	if !ws.IsMember(callerID) {
		return nil, shared.ErrForbidden
	}

	projects, err := s.projectRepo.FindByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	results := make([]*ProjectResponse, 0, len(projects))
	for _, p := range projects {
		results = append(results, toProjectResponse(p))
	}
	return results, nil
}

// Update modifies a project's details.
func (s *ProjectService) Update(
	ctx context.Context,
	id shared.ID,
	name, description, color string,
	viewType string,
) (*ProjectResponse, error) {
	p, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, shared.ErrNotFound
	}

	if err := p.Update(name, description, color, projectDomain.ViewType(viewType)); err != nil {
		return nil, err
	}

	if err := s.projectRepo.Save(ctx, p); err != nil {
		return nil, err
	}

	return toProjectResponse(p), nil
}

// Archive marks a project as archived.
func (s *ProjectService) Archive(ctx context.Context, id shared.ID) error {
	p, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return shared.ErrNotFound
	}

	if err := p.Archive(); err != nil {
		return err
	}

	return s.projectRepo.Save(ctx, p)
}

// Unarchive marks a project as not archived.
func (s *ProjectService) Unarchive(ctx context.Context, id shared.ID) error {
	p, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return shared.ErrNotFound
	}

	if err := p.Unarchive(); err != nil {
		return err
	}

	return s.projectRepo.Save(ctx, p)
}

// Delete removes a project by ID.
func (s *ProjectService) Delete(ctx context.Context, id shared.ID) error {
	p, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return shared.ErrNotFound
	}

	return s.projectRepo.Delete(ctx, id)
}

func toProjectResponse(p *projectDomain.Project) *ProjectResponse {
	return &ProjectResponse{
		ID:          p.ID.String(),
		WorkspaceID: p.WorkspaceID.String(),
		Name:        p.Name,
		Description: p.Description,
		Color:       p.Color,
		ViewType:    string(p.ViewType),
		IsArchived:  p.IsArchived,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}
