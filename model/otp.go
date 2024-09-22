package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/util"

)

type OTPPurpose string

const (
	OTPPurposeEmailVerification OTPPurpose = "email_verification"
	OTPPurposePasswordReset     OTPPurpose = "password_reset"
)

type OTPStatus string

const (
	OTPStatusUsed     OTPStatus = "used"
	OTPStatusActive   OTPStatus = "active"
	OTPStatusExpired  OTPStatus = "expired"
	OTPStatusInactive OTPStatus = "inactive"
)

type OTPChannel string

const (
	OTPChannelEmail OTPChannel = "email"
	OTPChannelSMS   OTPChannel = "sms"
)

type OTP struct {
	bun.BaseModel `bun:"otp,alias:o"`

	ID        string     `bun:"id,pk,type:varchar(21)"`
	OTP       string     `bun:"otp,notnull"`
	Purpose   OTPPurpose `bun:"purpose"`
	Status    OTPStatus  `bun:"status"`
	Channel   OTPChannel `bun:"channel"`
	Recipient string     `bun:"recipient,notnull"`
	ExpiresAt time.Time  `bun:"expires_at,notnull"`
	CreatedAt time.Time  `bun:"created_at,notnull"`
	UpdatedAt time.Time  `bun:"updated_at,notnull"`
}

func (o *OTP) SetTimestamps() {
	now := time.Now()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}
	o.UpdatedAt = now
}

func (o *OTP) SetId() {
	if o.ID == "" {
		o.ID = util.GenerateId()
	}
}

var _ bun.BeforeAppendModelHook = (*OTP)(nil)

func (o *OTP) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	o.SetTimestamps()
	o.SetId()
	return nil
}
