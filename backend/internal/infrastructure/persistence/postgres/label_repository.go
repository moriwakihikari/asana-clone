package postgres

import (
	"context"
	"errors"

	"asana-clone-backend/internal/domain/label"
	"asana-clone-backend/internal/domain/shared"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LabelRepository implements label.LabelRepository using PostgreSQL.
type LabelRepository struct {
	pool *pgxpool.Pool
}

var _ label.LabelRepository = (*LabelRepository)(nil)

func NewLabelRepository(pool *pgxpool.Pool) *LabelRepository {
	return &LabelRepository{pool: pool}
}

func (r *LabelRepository) FindByID(ctx context.Context, id shared.ID) (*label.Label, error) {
	query := `
		SELECT id, workspace_id, name, color, created_at
		FROM labels
		WHERE id = $1`

	l := &label.Label{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&l.ID, &l.WorkspaceID, &l.Name, &l.Color, &l.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return l, nil
}

func (r *LabelRepository) FindByWorkspaceID(ctx context.Context, workspaceID shared.ID) ([]*label.Label, error) {
	query := `
		SELECT id, workspace_id, name, color, created_at
		FROM labels
		WHERE workspace_id = $1
		ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []*label.Label
	for rows.Next() {
		l := &label.Label{}
		if err := rows.Scan(&l.ID, &l.WorkspaceID, &l.Name, &l.Color, &l.CreatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, rows.Err()
}

func (r *LabelRepository) Save(ctx context.Context, l *label.Label) error {
	query := `
		INSERT INTO labels (id, workspace_id, name, color, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			color = EXCLUDED.color`

	_, err := r.pool.Exec(ctx, query,
		l.ID, l.WorkspaceID, l.Name, l.Color, l.CreatedAt,
	)
	return err
}

func (r *LabelRepository) Delete(ctx context.Context, id shared.ID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM labels WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	return nil
}
