CREATE TABLE tasks (
       id UUID PRIMARY KEY,
       parent_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
       title VARCHAR(255) NOT NULL,
       description TEXT,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_parent_id ON tasks(parent_id);