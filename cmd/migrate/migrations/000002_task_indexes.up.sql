CREATE INDEX idx_tasks_ready
ON tasks (status, scheduled_at, priority);

CREATE INDEX idx_tasks_locked_at
ON tasks (locked_at);

CREATE INDEX idx_tasks_task_name
ON tasks (task_name);
