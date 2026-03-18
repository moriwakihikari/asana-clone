package postgres

import (
	"context"
	"errors"

	"asana-clone-backend/internal/domain/project"
	"asana-clone-backend/internal/domain/shared"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProjectRepository implements project.ProjectRepository using PostgreSQL.
type ProjectRepository struct {
	pool *pgxpool.Pool
}

var _ project.ProjectRepository = (*ProjectRepository)(nil)

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) FindByID(ctx context.Context, id shared.ID) (*project.Project, error) {
	query := `
		SELECT id, workspace_id, name, description, color, view_type, is_archived, created_at, updated_at
		FROM projects
		WHERE id = $1`

	p := &project.Project{}
	var viewType string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.WorkspaceID, &p.Name, &p.Description,
		&p.Color, &viewType, &p.IsArchived, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	p.ViewType = project.ViewType(viewType)
	return p, nil
}

func (r *ProjectRepository) FindByWorkspaceID(ctx context.Context, workspaceID shared.ID) ([]*project.Project, error) {
	query := `
		SELECT id, workspace_id, name, description, color, view_type, is_archived, created_at, updated_at
		FROM projects
		WHERE workspace_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*project.Project
	for rows.Next() {
		p := &project.Project{}
		var viewType string
		if err := rows.Scan(
			&p.ID, &p.WorkspaceID, &p.Name, &p.Description,
			&p.Color, &viewType, &p.IsArchived, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.ViewType = project.ViewType(viewType)
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) Save(ctx context.Context, p *project.Project) error {
	query := `
		INSERT INTO projects (id, workspace_id, name, description, color, view_type, is_archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			color = EXCLUDED.color,
			view_type = EXCLUDED.view_type,
			is_archived = EXCLUDED.is_archived,
			updated_at = EXCLUDED.updated_at`

	_, err := r.pool.Exec(ctx, query,
		p.ID, p.WorkspaceID, p.Name, p.Description,
		p.Color, string(p.ViewType), p.IsArchived, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id shared.ID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	return nil
}
