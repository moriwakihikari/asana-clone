package workspace

import (
	"context"

	"asana-clone-backend/internal/domain/shared"
)

type WorkspaceRepository interface {
	FindByID(ctx context.Context, id shared.ID) (*Workspace, error)
	FindByMemberID(ctx context.Context, userID shared.ID) ([]*Workspace, error)
	Save(ctx context.Context, workspace *Workspace) error
	Delete(ctx context.Context, id shared.ID) error
}
