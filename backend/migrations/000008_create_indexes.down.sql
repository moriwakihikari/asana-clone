-- comments indexes
DROP INDEX IF EXISTS idx_comments_task_created;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_task_id;

-- task_labels indexes
DROP INDEX IF EXISTS idx_task_labels_label_id;

-- tasks indexes
DROP INDEX IF EXISTS idx_tasks_search;
DROP INDEX IF EXISTS idx_tasks_assignee_active;
DROP INDEX IF EXISTS idx_tasks_active;
DROP INDEX IF EXISTS idx_tasks_section_position;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_assignee_id;
DROP INDEX IF EXISTS idx_tasks_section_id;
DROP INDEX IF EXISTS idx_tasks_project_id;

-- sections indexes
DROP INDEX IF EXISTS idx_sections_project_position;
DROP INDEX IF EXISTS idx_sections_project_id;

-- projects indexes
DROP INDEX IF EXISTS idx_projects_workspace_active;
DROP INDEX IF EXISTS idx_projects_workspace_id;

-- labels indexes
DROP INDEX IF EXISTS idx_labels_workspace_id;

-- workspace_members indexes
DROP INDEX IF EXISTS idx_workspace_members_user_id;
