package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/util"

)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID            string    `bun:"id,pk,type:varchar(21)"`
	Email         string    `bun:"email,unique,notnull"`
	EmailVerified bool      `bun:"email_verified,notnull,default:false"`
	Password      string    `bun:"password,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}

func (u *User) SetTimestamps() {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

func (u *User) SetId() {
	if u.ID == "" {
		u.ID = util.GenerateId()
	}
}

var _ bun.BeforeAppendModelHook = (*User)(nil)

func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	u.SetTimestamps()
	u.SetId()
	return nil
}
