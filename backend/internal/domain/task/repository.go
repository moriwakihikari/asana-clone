package task

import (
	"context"

	"asana-clone-backend/internal/domain/shared"
)

type TaskRepository interface {
	FindByID(ctx context.Context, id shared.ID) (*Task, error)
	FindByProjectID(ctx context.Context, projectID shared.ID) ([]*Task, error)
	FindWithFilters(ctx context.Context, filters TaskFilters) ([]*Task, error)
	Save(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id shared.ID) error
}
