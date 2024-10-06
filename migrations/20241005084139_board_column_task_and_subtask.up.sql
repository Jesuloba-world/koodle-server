SET
    statement_timeout = 0;

--bun:split
CREATE TABLE
    boards (
        id CHAR(21) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    columns (
        id CHAR(21) PRIMARY KEY,
        board_id CHAR(21) NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
        name VARCHAR(255) NOT NULL,
        position INTEGER NOT NULL,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    tasks (
        id CHAR(21) PRIMARY KEY,
        column_id CHAR(21) NOT NULL REFERENCES columns (id) ON DELETE CASCADE,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        position INTEGER NOT NULL,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    subtasks (
        id CHAR(21) PRIMARY KEY,
        task_id CHAR(21) NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
        name VARCHAR(255) NOT NULL,
        is_completed BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_columns_board_id ON columns (board_id);

CREATE INDEX idx_tasks_column_id ON tasks (column_id);

CREATE INDEX idx_subtasks_task_id ON subtasks (task_id);