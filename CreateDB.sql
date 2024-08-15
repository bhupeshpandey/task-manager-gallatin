CREATE TABLE tasks (
       id VARCHAR(255) PRIMARY KEY,
       parent_id VARCHAR(255),
       title VARCHAR(255) NOT NULL,
       description TEXT,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);