package project

import (
	"context"

	"asana-clone-backend/internal/domain/shared"
)

type ProjectRepository interface {
	FindByID(ctx context.Context, id shared.ID) (*Project, error)
	FindByWorkspaceID(ctx context.Context, workspaceID shared.ID) ([]*Project, error)
	Save(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id shared.ID) error
}
