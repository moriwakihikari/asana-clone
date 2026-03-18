package section

import (
	"context"

	"asana-clone-backend/internal/domain/shared"
)

type SectionRepository interface {
	FindByID(ctx context.Context, id shared.ID) (*Section, error)
	FindByProjectID(ctx context.Context, projectID shared.ID) ([]*Section, error)
	Save(ctx context.Context, section *Section) error
	Delete(ctx context.Context, id shared.ID) error
}
