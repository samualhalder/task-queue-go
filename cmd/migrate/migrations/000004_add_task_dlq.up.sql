CREATE TABLE IF NOT EXISTS tasks_dlq (
    id UUID PRIMARY KEY,
    task_id UUID,
    payload JSONB NOT NULL,
    attempts INT NOT NULL,
    error_text TEXT NOT NULL,
    error_type TEXT,
    failed_at TIMESTAMP NOT NULL DEFAULT NOW()
);
