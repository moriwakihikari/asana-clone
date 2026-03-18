CREATE TYPE project_view_type AS ENUM ('list', 'board', 'calendar', 'timeline');

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000),
    color CHAR(7) NOT NULL DEFAULT '#6366F1',
    view_type project_view_type NOT NULL DEFAULT 'board',
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
