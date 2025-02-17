package taskrepo

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/model"
)

type TaskRepo struct {
	db *bun.DB
}

func NewTaskRepo(db *bun.DB) *TaskRepo {
	return &TaskRepo{
		db: db,
	}
}

func (r *TaskRepo) AddTaskToColumn(ctx context.Context, task *model.Task, subtasks []*model.Subtask) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(task).Exec(ctx); err != nil {
			return fmt.Errorf("failed to insert task: %w", err)
		}

		if len(subtasks) > 0 {
			for _, subtask := range subtasks {
				subtask.TaskID = task.ID
			}
			if _, err := tx.NewInsert().Model(&subtasks).Exec(ctx); err != nil {
				return fmt.Errorf("failed to create columns: %w", err)
			}
		}

		task.Subtasks = make([]model.Subtask, len(subtasks))
		for i, st := range subtasks {
			task.Subtasks[i] = *st
		}

		return nil
	})
}
