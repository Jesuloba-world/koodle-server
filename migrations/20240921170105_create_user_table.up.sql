SET
    statement_timeout = 0;

--bun:split
CREATE TABLE
    users (
        id VARCHAR(21) PRIMARY KEY,
        email VARCHAR(255) NOT NULL UNIQUE,
        email_verified BOOLEAN NOT NULL DEFAULT FALSE,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

--bun:split
CREATE INDEX idx_users_email ON users (email);