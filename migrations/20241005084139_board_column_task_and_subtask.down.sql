SET
    statement_timeout = 0;

--bun:split
DROP INDEX IF EXISTS idx_subtasks_task_id;

DROP INDEX IF EXISTS idx_tasks_column_id;

DROP INDEX IF EXISTS idx_columns_board_id;

DROP TABLE IF EXISTS subtasks;

DROP TABLE IF EXISTS tasks;

DROP TABLE IF EXISTS columns;

DROP TABLE IF EXISTS boards;