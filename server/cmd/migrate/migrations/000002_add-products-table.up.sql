CREATE TABLE goals (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  priority VARCHAR(20) NOT NULL DEFAULT 'medium',
  status VARCHAR(20) NOT NULL DEFAULT 'todo',
  owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT goals_priority_check CHECK (priority IN ('high', 'medium', 'low')),
  CONSTRAINT goals_status_check CHECK (status IN ('todo', 'in_progress', 'achieved'))
);
