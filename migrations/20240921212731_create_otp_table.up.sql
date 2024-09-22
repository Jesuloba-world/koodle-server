SET
    statement_timeout = 0;

--bun:split
CREATE TABLE
    otp (
        id VARCHAR(21) PRIMARY KEY,
        otp VARCHAR(6) NOT NULL,
        recipient VARCHAR(255) NOT NULL,
        purpose VARCHAR(100) NOT NULL,
        status VARCHAR(20) NOT NULL,
        channel VARCHAR(20) NOT NULL,
        expires_at TIMESTAMP NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );