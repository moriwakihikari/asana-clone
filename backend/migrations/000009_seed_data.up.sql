-- Seed demo user
-- Password: "password123" hashed with bcrypt
INSERT INTO users (id, name, email, password_hash) VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Demo User',
    'demo@example.com',
    '$2a$10$fJVaquH.KwnfjWHWfbtRLedQ5LYc6pZvK6NrAJ3ZKFzelIiLyY5lC'
);

-- Seed workspace
INSERT INTO workspaces (id, name, description, owner_id) VALUES (
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'My Workspace',
    'Default workspace for demo user',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
);

-- Add demo user as workspace owner
INSERT INTO workspace_members (workspace_id, user_id, role) VALUES (
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'owner'
);

-- Seed project
INSERT INTO projects (id, workspace_id, name, description, color, view_type) VALUES (
    'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'Product Launch',
    'Track all tasks related to the product launch',
    '#6366F1',
    'board'
);

-- Seed sections
INSERT INTO sections (id, project_id, name, position) VALUES
    ('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'To Do',       1024),
    ('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'In Progress', 2048),
    ('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Done',        3072);

-- Seed tasks
INSERT INTO tasks (id, project_id, section_id, assignee_id, title, description, status, priority, due_date, position) VALUES
    (
        'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
        'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
        'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
        'Design landing page',
        'Create wireframes and high-fidelity mockups for the new landing page',
        'todo',
        'high',
        NOW() + INTERVAL '7 days',
        1024
    ),
    (
        'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
        'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
        'Write API documentation',
        'Document all REST API endpoints with request/response examples',
        'todo',
        'medium',
        NOW() + INTERVAL '14 days',
        2048
    ),
    (
        'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
        'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
        'Implement authentication flow',
        'Set up JWT-based authentication with refresh tokens',
        'in_progress',
        'urgent',
        NOW() + INTERVAL '3 days',
        1024
    ),
    (
        'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
        'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
        'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        NULL,
        'Set up CI/CD pipeline',
        'Configure GitHub Actions for automated testing and deployment',
        'in_progress',
        'low',
        NULL,
        2048
    ),
    (
        'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a05',
        'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
        'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
        'Set up project repository',
        'Initialize the monorepo with Go backend and React frontend',
        'done',
        'none',
        NULL,
        1024
    );
