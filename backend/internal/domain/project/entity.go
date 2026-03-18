package project

import (
	"strings"
	"time"

	"asana-clone-backend/internal/domain/shared"
)

type ViewType string

const (
	ViewTypeList     ViewType = "list"
	ViewTypeBoard    ViewType = "board"
	ViewTypeCalendar ViewType = "calendar"
	ViewTypeTimeline ViewType = "timeline"
)

func (v ViewType) IsValid() bool {
	switch v {
	case ViewTypeList, ViewTypeBoard, ViewTypeCalendar, ViewTypeTimeline:
		return true
	}
	return false
}

type Project struct {
	ID          shared.ID
	WorkspaceID shared.ID
	Name        string
	Description string
	Color       string
	ViewType    ViewType
	IsArchived  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProject(workspaceID shared.ID, name, description, color string, viewType ViewType) (*Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, shared.NewValidationError("INVALID_NAME", "project name must not be empty", "name")
	}
	if !viewType.IsValid() {
		return nil, shared.NewValidationError("INVALID_VIEW_TYPE", "invalid view type", "view_type")
	}

	now := time.Now()
	return &Project{
		ID:          shared.NewID(),
		WorkspaceID: workspaceID,
		Name:        name,
		Description: strings.TrimSpace(description),
		Color:       strings.TrimSpace(color),
		ViewType:    viewType,
		IsArchived:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (p *Project) Update(name, description, color string, viewType ViewType) error {
	if p.IsArchived {
		return shared.NewDomainError("PROJECT_ARCHIVED", "cannot update an archived project")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return shared.NewValidationError("INVALID_NAME", "project name must not be empty", "name")
	}
	if !viewType.IsValid() {
		return shared.NewValidationError("INVALID_VIEW_TYPE", "invalid view type", "view_type")
	}

	p.Name = name
	p.Description = strings.TrimSpace(description)
	p.Color = strings.TrimSpace(color)
	p.ViewType = viewType
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) Archive() error {
	if p.IsArchived {
		return shared.NewDomainError("ALREADY_ARCHIVED", "project is already archived")
	}
	p.IsArchived = true
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) Unarchive() error {
	if !p.IsArchived {
		return shared.NewDomainError("NOT_ARCHIVED", "project is not archived")
	}
	p.IsArchived = false
	p.UpdatedAt = time.Now()
	return nil
}
