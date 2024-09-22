SET
    statement_timeout = 0;

--bun:split
DROP INDEX IF EXISTS idx_users_email;

--bun:split
DROP TABLE IF EXISTS users;