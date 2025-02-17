package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/util"

)

type Task struct {
	bun.BaseModel `bun:"table:tasks,alias:t"`

	ID          string    `bun:"id,pk,type:char(21)" json:"id"`
	ColumnID    string    `bun:"column_id,notnull,type:char(21)" json:"column_id"`
	Title       string    `bun:"title,notnull" json:"title"`
	Description string    `bun:"description" json:"description"`
	Position    int       `bun:"position,notnull" json:"position"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	Subtasks []Subtask `bun:"rel:has-many,join:id=task_id" json:"subtasks"`
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
