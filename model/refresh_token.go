package model

import (
	"time"

	"github.com/uptrace/bun"
)

type RefreshToken struct {
	bun.BaseModel `bun:"table:refresh_tokens,alias:r"`

	RefreshToken string    `bun:"refresh_token,pk"`
	UserID       string    `bun:"user_id"`
	ExpiresAt    time.Time `bun:"expires_at"`
	Revoked      bool      `bun:"revoked,default:false"`
}
