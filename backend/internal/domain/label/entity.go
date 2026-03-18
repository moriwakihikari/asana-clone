package label

import (
	"strings"
	"time"

	"asana-clone-backend/internal/domain/shared"
)

type Label struct {
	ID          shared.ID
	WorkspaceID shared.ID
	Name        string
	Color       string
	CreatedAt   time.Time
}

func NewLabel(workspaceID shared.ID, name, color string) (*Label, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, shared.NewValidationError("INVALID_NAME", "label name must not be empty", "name")
	}

	color = strings.TrimSpace(color)
	if color == "" {
		return nil, shared.NewValidationError("INVALID_COLOR", "label color must not be empty", "color")
	}

	return &Label{
		ID:          shared.NewID(),
		WorkspaceID: workspaceID,
		Name:        name,
		Color:       color,
		CreatedAt:   time.Now(),
	}, nil
}
