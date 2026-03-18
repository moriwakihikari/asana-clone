package postgres

import (
	"context"
	"errors"

	"asana-clone-backend/internal/domain/shared"
	"asana-clone-backend/internal/domain/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository implements user.UserRepository using PostgreSQL.
type UserRepository struct {
	pool *pgxpool.Pool
}

// Compile-time check that UserRepository implements user.UserRepository.
var _ user.UserRepository = (*UserRepository)(nil)

// NewUserRepository creates a new UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) FindByID(ctx context.Context, id shared.ID) (*user.User, error) {
	query := `
		SELECT id, name, email, password_hash, COALESCE(avatar_url, ''), created_at, updated_at
		FROM users
		WHERE id = $1`

	u := &user.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, name, email, password_hash, COALESCE(avatar_url, ''), created_at, updated_at
		FROM users
		WHERE email = $1`

	u := &user.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = EXCLUDED.updated_at`

	_, err := r.pool.Exec(ctx, query,
		u.ID, u.Name, u.Email, u.PasswordHash,
		u.AvatarURL, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id shared.ID) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	return nil
}
