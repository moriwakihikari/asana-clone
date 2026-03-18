package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"asana-clone-backend/internal/domain/shared"
	"asana-clone-backend/internal/domain/task"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TaskRepository implements task.TaskRepository using PostgreSQL.
type TaskRepository struct {
	pool *pgxpool.Pool
}

var _ task.TaskRepository = (*TaskRepository)(nil)

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) FindByID(ctx context.Context, id shared.ID) (*task.Task, error) {
	query := `
		SELECT id, project_id, section_id, assignee_id, title, description,
		       status, priority, due_date, position, created_at, updated_at
		FROM tasks
		WHERE id = $1`

	t := &task.Task{}
	var status, priority string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.ProjectID, &t.SectionID, &t.AssigneeID,
		&t.Title, &t.Description, &status, &priority,
		&t.DueDate, &t.Position, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	t.Status = task.Status(status)
	t.Priority = task.Priority(priority)

	labelIDs, err := r.findTaskLabelIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	t.LabelIDs = labelIDs

	return t, nil
}

func (r *TaskRepository) FindByProjectID(ctx context.Context, projectID shared.ID) ([]*task.Task, error) {
	query := `
		SELECT id, project_id, section_id, assignee_id, title, description,
		       status, priority, due_date, position, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
		ORDER BY position ASC`

	return r.queryTasks(ctx, query, projectID)
}

func (r *TaskRepository) FindWithFilters(ctx context.Context, filters task.TaskFilters) ([]*task.Task, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filters.ProjectID != nil {
		conditions = append(conditions, fmt.Sprintf("t.project_id = $%d", argIdx))
		args = append(args, *filters.ProjectID)
		argIdx++
	}
	if filters.SectionID != nil {
		conditions = append(conditions, fmt.Sprintf("t.section_id = $%d", argIdx))
		args = append(args, *filters.SectionID)
		argIdx++
	}
	if filters.AssigneeID != nil {
		conditions = append(conditions, fmt.Sprintf("t.assignee_id = $%d", argIdx))
		args = append(args, *filters.AssigneeID)
		argIdx++
	}
	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("t.status = $%d", argIdx))
		args = append(args, string(*filters.Status))
		argIdx++
	}
	if filters.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("t.priority = $%d", argIdx))
		args = append(args, string(*filters.Priority))
		argIdx++
	}
	if filters.DueBefore != nil {
		conditions = append(conditions, fmt.Sprintf("t.due_date <= $%d", argIdx))
		args = append(args, *filters.DueBefore)
		argIdx++
	}
	if filters.DueAfter != nil {
		conditions = append(conditions, fmt.Sprintf("t.due_date >= $%d", argIdx))
		args = append(args, *filters.DueAfter)
		argIdx++
	}
	if filters.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(t.title ILIKE $%d OR t.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filters.Query+"%")
		argIdx++
	}
	if filters.LabelID != nil {
		conditions = append(conditions, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM task_labels tl WHERE tl.task_id = t.id AND tl.label_id = $%d)", argIdx))
		args = append(args, *filters.LabelID)
		argIdx++
	}

	query := `
		SELECT t.id, t.project_id, t.section_id, t.assignee_id, t.title, t.description,
		       t.status, t.priority, t.due_date, t.position, t.created_at, t.updated_at
		FROM tasks t`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY t.position ASC"

	return r.queryTasks(ctx, query, args...)
}

func (r *TaskRepository) Save(ctx context.Context, t *task.Task) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO tasks (id, project_id, section_id, assignee_id, title, description,
		                   status, priority, due_date, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			project_id = EXCLUDED.project_id,
			section_id = EXCLUDED.section_id,
			assignee_id = EXCLUDED.assignee_id,
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			status = EXCLUDED.status,
			priority = EXCLUDED.priority,
			due_date = EXCLUDED.due_date,
			position = EXCLUDED.position,
			updated_at = EXCLUDED.updated_at`

	_, err = tx.Exec(ctx, query,
		t.ID, t.ProjectID, t.SectionID, t.AssigneeID,
		t.Title, t.Description, string(t.Status), string(t.Priority),
		t.DueDate, t.Position, t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Sync task_labels: delete all, then re-insert.
	_, err = tx.Exec(ctx, `DELETE FROM task_labels WHERE task_id = $1`, t.ID)
	if err != nil {
		return err
	}

	if len(t.LabelIDs) > 0 {
		for _, labelID := range t.LabelIDs {
			_, err = tx.Exec(ctx,
				`INSERT INTO task_labels (task_id, label_id) VALUES ($1, $2)`,
				t.ID, labelID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *TaskRepository) Delete(ctx context.Context, id shared.ID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM task_labels WHERE task_id = $1`, id)
	if err != nil {
		return err
	}

	result, err := tx.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return shared.ErrNotFound
	}

	return tx.Commit(ctx)
}

// queryTasks is a helper that scans task rows and loads their labels.
func (r *TaskRepository) queryTasks(ctx context.Context, query string, args ...interface{}) ([]*task.Task, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		t := &task.Task{}
		var status, priority string
		if err := rows.Scan(
			&t.ID, &t.ProjectID, &t.SectionID, &t.AssigneeID,
			&t.Title, &t.Description, &status, &priority,
			&t.DueDate, &t.Position, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		t.Status = task.Status(status)
		t.Priority = task.Priority(priority)
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load labels for each task.
	for _, t := range tasks {
		labelIDs, err := r.findTaskLabelIDs(ctx, t.ID)
		if err != nil {
			return nil, err
		}
		t.LabelIDs = labelIDs
	}

	return tasks, nil
}

func (r *TaskRepository) findTaskLabelIDs(ctx context.Context, taskID shared.ID) ([]shared.ID, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT label_id FROM task_labels WHERE task_id = $1`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []shared.ID
	for rows.Next() {
		var id shared.ID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = make([]shared.ID, 0)
	}
	return ids, rows.Err()
}
