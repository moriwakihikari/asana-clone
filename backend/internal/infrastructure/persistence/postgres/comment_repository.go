package postgres

import (
	"context"
	"errors"

	"asana-clone-backend/internal/domain/comment"
	"asana-clone-backend/internal/domain/shared"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CommentRepository implements comment.CommentRepository using PostgreSQL.
type CommentRepository struct {
	pool *pgxpool.Pool
}

var _ comment.CommentRepository = (*CommentRepository)(nil)

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

func (r *CommentRepository) FindByID(ctx context.Context, id shared.ID) (*comment.Comment, error) {
	query := `
		SELECT id, task_id, user_id, content, is_edited, created_at, updated_at
		FROM comments
		WHERE id = $1`

	c := &comment.Comment{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.TaskID, &c.UserID, &c.Content,
		&c.IsEdited, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *CommentRepository) FindByTaskID(ctx context.Context, taskID shared.ID) ([]*comment.Comment, error) {
	query := `
		SELECT id, task_id, user_id, content, is_edited, created_at, updated_at
		FROM comments
		WHERE task_id = $1
		ORDER BY created_at ASC`

	rows, err := r.pool.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		c := &comment.Comment{}
		if err := rows.Scan(
			&c.ID, &c.TaskID, &c.UserID, &c.Content,
			&c.IsEdited, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *CommentRepository) Save(ctx context.Context, c *comment.Comment) error {
	query := `
		INSERT INTO comments (id, task_id, user_id, content, is_edited, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			is_edited = EXCLUDED.is_edited,
			updated_at = EXCLUDED.updated_at`

	_, err := r.pool.Exec(ctx, query,
		c.ID, c.TaskID, c.UserID, c.Content,
		c.IsEdited, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CommentRepository) Delete(ctx context.Context, id shared.ID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return shared.ErrNotFound
	}
	return nil
}
