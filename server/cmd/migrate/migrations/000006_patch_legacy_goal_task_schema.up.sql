-- Compatibility migration for DBs created from older migration definitions.
-- Adds columns expected by current application code if they are missing.

ALTER TABLE goals
  ADD COLUMN IF NOT EXISTS priority VARCHAR(20),
  ADD COLUMN IF NOT EXISTS status VARCHAR(20);

UPDATE goals SET priority = 'medium' WHERE priority IS NULL;
UPDATE goals SET status = 'todo' WHERE status IS NULL;

ALTER TABLE goals ALTER COLUMN priority SET DEFAULT 'medium';
ALTER TABLE goals ALTER COLUMN priority SET NOT NULL;
ALTER TABLE goals ALTER COLUMN status SET DEFAULT 'todo';
ALTER TABLE goals ALTER COLUMN status SET NOT NULL;

ALTER TABLE tasks
  ADD COLUMN IF NOT EXISTS priority VARCHAR(20),
  ADD COLUMN IF NOT EXISTS is_completed BOOLEAN;

UPDATE tasks SET priority = 'medium' WHERE priority IS NULL;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'tasks'
      AND column_name = 'status'
  ) THEN
    UPDATE tasks
    SET is_completed = CASE WHEN status = 'done' THEN TRUE ELSE FALSE END
    WHERE is_completed IS NULL;
  ELSE
    UPDATE tasks
    SET is_completed = FALSE
    WHERE is_completed IS NULL;
  END IF;
END $$;

ALTER TABLE tasks ALTER COLUMN priority SET DEFAULT 'medium';
ALTER TABLE tasks ALTER COLUMN priority SET NOT NULL;
ALTER TABLE tasks ALTER COLUMN is_completed SET DEFAULT FALSE;
ALTER TABLE tasks ALTER COLUMN is_completed SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'goals_priority_check'
  ) THEN
    ALTER TABLE goals
      ADD CONSTRAINT goals_priority_check
      CHECK (priority IN ('high', 'medium', 'low'));
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'goals_status_check'
  ) THEN
    ALTER TABLE goals
      ADD CONSTRAINT goals_status_check
      CHECK (status IN ('todo', 'in_progress', 'achieved'));
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'tasks_priority_check'
  ) THEN
    ALTER TABLE tasks
      ADD CONSTRAINT tasks_priority_check
      CHECK (priority IN ('high', 'medium', 'low'));
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_goals_status_priority ON goals(status, priority);
CREATE INDEX IF NOT EXISTS idx_tasks_completion_priority ON tasks(is_completed, priority);
