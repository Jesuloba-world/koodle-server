package otpService

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/uptrace/bun"

	"github.com/Jesuloba-world/koodle-server/model"
)

type Sender interface {
	SendOTPToRecipient(channel model.OTPChannel, purpose model.OTPPurpose, recipient string, otp string, expiration time.Duration) error
}

type OTPService struct {
	db                    *bun.DB
	otpExpirationDuration time.Duration
	otpGenerateTimeLapse  time.Duration
	senderService         Sender
}

func NewOTPService(
	db *bun.DB,
	otpExpirationDuration,
	otpGenerateTimeLapse time.Duration,
	senderService Sender,
) *OTPService {
	return &OTPService{
		db:                    db,
		otpExpirationDuration: otpExpirationDuration,
		otpGenerateTimeLapse:  otpGenerateTimeLapse,
		senderService:         senderService,
	}
}

func (s *OTPService) generateOTP(length ...int) string {
	// Default OTP length is 4
	otpLength := 6
	if len(length) > 0 && length[0] > 0 {
		otpLength = length[0]
	}

	min := int(math.Pow10(otpLength - 1))
	max := int(math.Pow10(otpLength)) - 1

	otp := rand.Intn(max-min+1) + min

	return fmt.Sprintf("%0*d", otpLength, otp)
}

func (s *OTPService) GenerateAndStoreOTP(purpose model.OTPPurpose, channel model.OTPChannel, recipient string) (string, error) {
	otp := s.generateOTP()

	newOtp := &model.OTP{
		OTP:       otp,
		Purpose:   purpose,
		Channel:   channel,
		Recipient: recipient,
		Status:    model.OTPStatusActive,
		ExpiresAt: time.Now().Add(s.otpExpirationDuration),
	}

	_, err := s.db.NewInsert().Model(newOtp).Exec(context.Background())
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (s *OTPService) VerifyOTP(purpose model.OTPPurpose, channel model.OTPChannel, recipient string, otp string, markUsed bool) (bool, error) {
	//  if markUsed is true, the otp will be marked as used after verification
	//  if markUsed is false, the otp will not be marked as used after verification, only validated

	var storedOTP model.OTP

	err := s.db.NewSelect().Model(&storedOTP).
		Where("purpose = ?", purpose).
		Where("channel = ?", channel).
		Where("recipient = ?", recipient).
		Where("status = ?", model.OTPStatusActive).
		Where("expires_at > ?", time.Now()).
		Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if storedOTP.OTP != otp {
		return false, nil
	}

	if markUsed {
		storedOTP.Status = model.OTPStatusUsed
		_, err := s.db.NewUpdate().Model(&storedOTP).Where("id = ?", storedOTP.ID).Exec(context.Background())
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *OTPService) ThrottleOTP(recipient string, purpose model.OTPPurpose) (bool, error) {
	var lastOTPCreatedAt time.Time
	err := s.db.NewSelect().
		Model((*model.OTP)(nil)).
		Column("created_at").
		Where("recipient = ?", recipient).
		Where("purpose = ?", purpose).
		Where("status IN (?)", bun.In([]model.OTPStatus{model.OTPStatusActive, model.OTPStatusExpired})).
		Order("created_at DESC").
		Limit(1).
		Scan(context.Background(), &lastOTPCreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		// no active or expired OTP found, allow only one to be generated at a time
		return true, nil
	}

	if err != nil {
		return false, err
	}

	if time.Since(lastOTPCreatedAt) < s.otpGenerateTimeLapse {
		// last otp was generated within timelapse, so don't generate a new one
		return false, nil
	}

	// allow generating new OTP
	return true, nil
}

func (s *OTPService) SendOTP(purpose model.OTPPurpose, channel model.OTPChannel, recipient string) (string, error) {
	// deactivate any active OTPs for the same recipient and purpose
	_, err := s.db.NewUpdate().
		Model((*model.OTP)(nil)).
		Set("status = ?", model.OTPStatusInactive).
		Where("recipient = ?", recipient).
		Where("purpose = ?", purpose).
		Where("status = ?", model.OTPStatusActive).
		Exec(context.Background())

	if err != nil {
		return "", err
	}

	otp, err := s.GenerateAndStoreOTP(purpose, channel, recipient)
	if err != nil {
		return "", nil
	}

	err = s.senderService.SendOTPToRecipient(channel, purpose, recipient, otp, s.otpExpirationDuration)
	if err != nil {
		return "", err
	}

	return otp, nil
}
