package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/util"

)

type Task struct {
	bun.BaseModel `bun:"table:tasks,alias:t"`

	ID          string    `bun:"id,pk,type:char(21)"`
	ColumnID    string    `bun:"column_id,notnull,type:char(21)"`
	Title       string    `bun:"title,notnull"`
	Description string    `bun:"description"`
	Position    int       `bun:"position,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp"`

	Subtasks []Subtask `bun:"rel:has-many,join:id=task_id"`
}

func (u *Task) SetTimestamps() {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

func (u *Task) SetId() {
	if u.ID == "" {
		u.ID = util.GenerateId()
	}
}

var _ bun.BeforeAppendModelHook = (*Task)(nil)

func (u *Task) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	u.SetTimestamps()
	u.SetId()
	return nil
}
