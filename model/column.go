package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/util"
)

type Column struct {
	bun.BaseModel `bun:"table:columns,alias:c"`

	ID        string    `bun:"id,pk,type:char(21)" json:"id"`
	BoardID   string    `bun:"board_id,notnull,type:char(21)" json:"-"`
	Name      string    `bun:"name,notnull" json:"name"`
	Position  int       `bun:"position,notnull" json:"position"`
	Color     string    `bun:"color,notnull" json:"color"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`

	Tasks []Task `bun:"rel:has-many,join:id=column_id" json:"tasks"`
}

func (u *Column) SetTimestamps() {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

func (u *Column) SetId() {
	if u.ID == "" {
		u.ID = util.GenerateId()
	}
}

func (u *Column) SetColor() {
	if u.Color == "" {
		u.Color = util.GenerateColor()
	}
}

var _ bun.BeforeAppendModelHook = (*Column)(nil)

func (u *Column) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	u.SetTimestamps()
	u.SetId()
	u.SetColor()
	return nil
}
