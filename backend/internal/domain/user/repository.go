package user

import (
	"context"

	"asana-clone-backend/internal/domain/shared"
)

type UserRepository interface {
	FindByID(ctx context.Context, id shared.ID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user *User) error
	Delete(ctx context.Context, id shared.ID) error
}
