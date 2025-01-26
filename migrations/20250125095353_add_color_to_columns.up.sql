SET statement_timeout = 0;

--bun:split

ALTER TABLE columns
ADD COLUMN color CHAR(7) NOT NULL DEFAULT '#8471F2';

