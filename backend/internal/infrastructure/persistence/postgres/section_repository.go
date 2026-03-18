package postgres

import (
	"context"
	"errors"

	"asana-clone-backend/internal/domain/section"
	"asana-clone-backend/internal/domain/shared"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SectionRepository implements section.SectionRepository using PostgreSQL.
type SectionRepository struct {
	pool *pgxpool.Pool
}

var _ section.SectionRepository = (*SectionRepository)(nil)

func NewSectionRepository(pool *pgxpool.Pool) *SectionRepository {
	return &SectionRepository{pool: pool}
}

func (r *SectionRepository) FindByID(ctx context.Context, id shared.ID) (*section.Section, error) {
	query := `
		SELECT id, project_id, name, position, created_at
		FROM sections
		WHERE id = $1`

	s := &section.Section{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.ProjectID, &s.Name, &s.Position, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SectionRepository) FindByProjectID(ctx context.Context, projectID shared.ID) ([]*section.Section, error) {
	query := `
		SELECT id, project_id, name, position, created_at
		FROM sections
		WHERE project_id = $1
		ORDER BY position ASC`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []*section.Section
	for rows.Next() {
		s := &section.Section{}
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.Position, &s.CreatedAt); err != nil {
			return nil, err
		}
		sections = append(sections, s)
	}
	return sections, rows.Err()
}

func (r *SectionRepository) Save(ctx context.Context, s *section.Section) error {
	query := `
		INSERT INTO sections (id, project_id, name, position, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			position = EXCLUDED.position`

	_, err := r.pool.Exec(ctx, query,
		s.ID, s.ProjectID, s.Name, s.Position, s.CreatedAt,
	)
	return err
}

func (r *SectionRepository) Delete(ctx context.Context, id shared.ID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM sections WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	return nil
}
