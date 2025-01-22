SET statement_timeout = 0;

--bun:split

DROP INDEX IF EXISTS idx_boards_user_id

--bun:split

ALTER TABLE boards
DROP COLUMN user_id
