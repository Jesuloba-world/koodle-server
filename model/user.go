package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	"github.com/Jesuloba-world/koodle-server/util"

)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID            string    `bun:"id,pk,type:varchar(21)"`
	Email         string    `bun:"email,unique,notnull"`
	EmailVerified bool      `bun:"email_verified,notnull,default:false"`
	Password      string    `bun:"password,notnull" json:"-"`
	CreatedAt     time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,default:current_timestamp"`

	Boards []Board `bun:"rel:has-many,join:id=user_id"`
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

func (u *User) SetPassword(password string) error {
	// hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

var _ bun.BeforeAppendModelHook = (*User)(nil)

func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	u.SetTimestamps()
	u.SetId()
	return nil
}

type UserResponse struct {
	ID            string    `json:"id" example:"136789874673893" doc:"ID of user"`
	Email         string    `json:"email" example:"user@example.com" doc:"Email address of user"`
	EmailVerified bool      `json:"email_verified" example:"true" doc:"Whether user's email is verified"`
	CreatedAt     time.Time `json:"created_at" doc:"Time user was created"`
	UpdatedAt     time.Time `json:"updated_at" doc:"Time user was last updated"`
}

func (u *User) MapUserToResponse() *UserResponse {
	return &UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
