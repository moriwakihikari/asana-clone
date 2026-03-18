package comment

import (
	"context"

	"asana-clone-backend/internal/domain/shared"
)

type CommentRepository interface {
	FindByID(ctx context.Context, id shared.ID) (*Comment, error)
	FindByTaskID(ctx context.Context, taskID shared.ID) ([]*Comment, error)
	Save(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id shared.ID) error
}
