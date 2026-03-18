package postgres

import (
	"context"
	"errors"

	"asana-clone-backend/internal/domain/shared"
	"asana-clone-backend/internal/domain/workspace"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WorkspaceRepository implements workspace.WorkspaceRepository using PostgreSQL.
type WorkspaceRepository struct {
	pool *pgxpool.Pool
}

var _ workspace.WorkspaceRepository = (*WorkspaceRepository)(nil)

func NewWorkspaceRepository(pool *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{pool: pool}
}

func (r *WorkspaceRepository) FindByID(ctx context.Context, id shared.ID) (*workspace.Workspace, error) {
	query := `
		SELECT id, name, description, owner_id, created_at
		FROM workspaces
		WHERE id = $1`

	ws := &workspace.Workspace{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&ws.ID, &ws.Name, &ws.Description, &ws.OwnerID, &ws.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	members, err := r.findMembers(ctx, id)
	if err != nil {
		return nil, err
	}
	ws.Members = members

	return ws, nil
}

func (r *WorkspaceRepository) FindByMemberID(ctx context.Context, userID shared.ID) ([]*workspace.Workspace, error) {
	query := `
		SELECT w.id, w.name, w.description, w.owner_id, w.created_at
		FROM workspaces w
		INNER JOIN workspace_members wm ON w.id = wm.workspace_id
		WHERE wm.user_id = $1
		ORDER BY w.created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*workspace.Workspace
	for rows.Next() {
		ws := &workspace.Workspace{}
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.Description, &ws.OwnerID, &ws.CreatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, ws)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load members for each workspace.
	for _, ws := range workspaces {
		members, err := r.findMembers(ctx, ws.ID)
		if err != nil {
			return nil, err
		}
		ws.Members = members
	}

	return workspaces, nil
}

func (r *WorkspaceRepository) Save(ctx context.Context, ws *workspace.Workspace) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Upsert workspace.
	wsQuery := `
		INSERT INTO workspaces (id, name, description, owner_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			owner_id = EXCLUDED.owner_id`

	_, err = tx.Exec(ctx, wsQuery,
		ws.ID, ws.Name, ws.Description, ws.OwnerID, ws.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Delete existing members and re-insert (simplest approach for consistency).
	_, err = tx.Exec(ctx, `DELETE FROM workspace_members WHERE workspace_id = $1`, ws.ID)
	if err != nil {
		return err
	}

	memberQuery := `
		INSERT INTO workspace_members (workspace_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)`

	for _, m := range ws.Members {
		_, err = tx.Exec(ctx, memberQuery, m.WorkspaceID, m.UserID, string(m.Role), m.JoinedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *WorkspaceRepository) Delete(ctx context.Context, id shared.ID) error {
	// workspace_members should be cascade-deleted via FK, but we delete explicitly for safety.
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM workspace_members WHERE workspace_id = $1`, id)
	if err != nil {
		return err
	}

	result, err := tx.Exec(ctx, `DELETE FROM workspaces WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return shared.ErrNotFound
	}

	return tx.Commit(ctx)
}

func (r *WorkspaceRepository) findMembers(ctx context.Context, workspaceID shared.ID) ([]workspace.WorkspaceMember, error) {
	query := `
		SELECT workspace_id, user_id, role, joined_at
		FROM workspace_members
		WHERE workspace_id = $1
		ORDER BY joined_at ASC`

	rows, err := r.pool.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []workspace.WorkspaceMember
	for rows.Next() {
		var m workspace.WorkspaceMember
		var role string
		if err := rows.Scan(&m.WorkspaceID, &m.UserID, &role, &m.JoinedAt); err != nil {
			return nil, err
		}
		m.Role = workspace.Role(role)
		members = append(members, m)
	}
	return members, rows.Err()
}
