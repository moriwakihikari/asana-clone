-- Remove seed data in reverse order of insertion
DELETE FROM tasks WHERE id IN (
    'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
    'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
    'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
    'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
    'e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a05'
);

DELETE FROM sections WHERE id IN (
    'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
    'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
    'd3eebc99-9c0b-4ef8-bb6d-6bb9bd380a03'
);

DELETE FROM projects WHERE id = 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33';

DELETE FROM workspace_members WHERE workspace_id = 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22';

DELETE FROM workspaces WHERE id = 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22';

DELETE FROM users WHERE id = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';
