-- Initialize database for Task Manager
-- This file is executed when PostgreSQL container starts for the first time

-- Connect to the task_db database
\c task_db;

-- Create tasks table (GORM will handle this, but keeping for reference)
-- The application uses GORM auto-migration, so this is optional

-- Grant additional permissions if needed
GRANT ALL PRIVILEGES ON DATABASE task_db TO taskuser;

-- Insert sample data (optional)
-- This can be uncommented if you want some initial test data

/*
INSERT INTO tasks (uuid, title, description, status, user_id, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Sample Task 1', 'This is a sample task for testing', 'Pending', 1, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440001', 'Sample Task 2', 'Another sample task', 'InProgress', 2, NOW(), NOW());
*/
