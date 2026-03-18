-- workspace_members indexes
CREATE INDEX idx_workspace_members_user_id ON workspace_members(user_id);

-- labels indexes
CREATE INDEX idx_labels_workspace_id ON labels(workspace_id);

-- projects indexes
CREATE INDEX idx_projects_workspace_id ON projects(workspace_id);
CREATE INDEX idx_projects_workspace_active ON projects(workspace_id) WHERE is_archived = FALSE;

-- sections indexes
CREATE INDEX idx_sections_project_id ON sections(project_id);
CREATE INDEX idx_sections_project_position ON sections(project_id, position);

-- tasks indexes
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_section_id ON tasks(section_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_due_date ON tasks(due_date) WHERE due_date IS NOT NULL;
CREATE INDEX idx_tasks_section_position ON tasks(section_id, position);
CREATE INDEX idx_tasks_active ON tasks(project_id, status) WHERE status NOT IN ('done', 'cancelled');
CREATE INDEX idx_tasks_assignee_active ON tasks(assignee_id, status) WHERE assignee_id IS NOT NULL AND status NOT IN ('done', 'cancelled');

-- GIN full-text search index on tasks
CREATE INDEX idx_tasks_search ON tasks USING GIN (to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- task_labels indexes
CREATE INDEX idx_task_labels_label_id ON task_labels(label_id);

-- comments indexes
CREATE INDEX idx_comments_task_id ON comments(task_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_task_created ON comments(task_id, created_at);
