SET statement_timeout = 0;

--bun:split

ALTER TABLE boards
ADD COLUMN user_id CHAR(21) NOT NULL REFERENCES users(id) ON DELETE CASCADE

--bun:split

CREATE INDEX idx_boards_user_id ON boards(user_id)
