package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/util"
)

type Subtask struct {
	bun.BaseModel `bun:"table:subtasks,alias:st"`

	ID          string    `bun:"id,pk,type:char(21)"`
	TaskID      string    `bun:"task_id,notnull,type:char(21)"`
	Name        string    `bun:"name,notnull"`
	IsCompleted bool      `bun:"is_completed,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

func (u *Subtask) SetTimestamps() {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

func (u *Subtask) SetId() {
	if u.ID == "" {
		u.ID = util.GenerateId()
	}
}

var _ bun.BeforeAppendModelHook = (*Subtask)(nil)

func (u *Subtask) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	u.SetTimestamps()
	u.SetId()
	return nil
}
