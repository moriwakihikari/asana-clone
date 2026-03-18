package label

import (
	"context"

	"asana-clone-backend/internal/domain/shared"
)

type LabelRepository interface {
	FindByID(ctx context.Context, id shared.ID) (*Label, error)
	FindByWorkspaceID(ctx context.Context, workspaceID shared.ID) ([]*Label, error)
	Save(ctx context.Context, label *Label) error
	Delete(ctx context.Context, id shared.ID) error
}
