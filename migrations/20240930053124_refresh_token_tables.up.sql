SET
    statement_timeout = 0;

--bun:split
CREATE TABLE
    "refresh_tokens" (
        refresh_token VARCHAR(255),
        user_id VARCHAR(255) NOT NULL,
        expires_at TIMESTAMP NOT NULL,
        revoked BOOLEAN NOT NULL DEFAULT FALSE
    );