-- Create the tasks table based on the Task model
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    uuid CHAR(36) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL,
    user_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance based on actual query patterns

-- 1. UUID lookup (most critical - keep as unique)
CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_uuid ON tasks(uuid);

-- 2. Composite index for List() method - handles all filter combinations + ordering
-- This single index covers: status, user_id, priority filters + created_at ordering
CREATE INDEX IF NOT EXISTS idx_tasks_composite_list ON tasks(status, user_id, priority, created_at DESC);

-- 3. Composite index for duplicate check (ExistsByTitleAndUser method)
CREATE INDEX IF NOT EXISTS idx_tasks_title_user_id ON tasks(title, user_id);

-- 4. Fallback indexes for partial filtering scenarios
CREATE INDEX IF NOT EXISTS idx_tasks_user_created ON tasks(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_status_created ON tasks(status, created_at DESC);
