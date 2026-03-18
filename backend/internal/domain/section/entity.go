package section

import (
	"strings"
	"time"

	"asana-clone-backend/internal/domain/shared"
)

type Section struct {
	ID        shared.ID
	ProjectID shared.ID
	Name      string
	Position  int
	CreatedAt time.Time
}

func NewSection(projectID shared.ID, name string, position int) (*Section, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, shared.NewValidationError("INVALID_NAME", "section name must not be empty", "name")
	}
	if position < 0 {
		return nil, shared.NewValidationError("INVALID_POSITION", "position must be non-negative", "position")
	}

	return &Section{
		ID:        shared.NewID(),
		ProjectID: projectID,
		Name:      name,
		Position:  position,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Section) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return shared.NewValidationError("INVALID_NAME", "section name must not be empty", "name")
	}
	s.Name = name
	return nil
}

func (s *Section) MoveTo(position int) error {
	if position < 0 {
		return shared.NewValidationError("INVALID_POSITION", "position must be non-negative", "position")
	}
	s.Position = position
	return nil
}
