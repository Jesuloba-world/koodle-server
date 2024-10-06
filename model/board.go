package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/util"

)

type Board struct {
	bun.BaseModel `bun:"table:boards,alias:b"`

	ID        string    `bun:"id,pk,type:char(21)" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`

	Columns []Column `bun:"rel:has-many,join:id=board_id" json:"columns"`
}

func (u *Board) SetTimestamps() {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

func (u *Board) SetId() {
	if u.ID == "" {
		u.ID = util.GenerateId()
	}
}

var _ bun.BeforeAppendModelHook = (*Board)(nil)

func (u *Board) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	u.SetTimestamps()
	u.SetId()
	return nil
}
