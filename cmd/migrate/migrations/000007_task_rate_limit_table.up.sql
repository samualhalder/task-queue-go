CREATE TABLE task_rate_limit (
    key TEXT PRIMARY KEY,

    capacity INT NOT NULL CHECK (capacity > 0),
    refill_rate NUMERIC NOT NULL CHECK (refill_rate > 0),

    tokens NUMERIC NOT NULL CHECK (tokens >= 0),
    last_refill_at TIMESTAMPTZ NOT NULL
);
