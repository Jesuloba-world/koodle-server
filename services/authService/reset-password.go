package authservice

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Jesuloba-world/koodle-server/lib/validator"
	"github.com/Jesuloba-world/koodle-server/model"
)

// start reset password
// resend password reset otp
// verify password reset otp
// reset password

func (s *AuthService) startResetPassword(ctx context.Context, req *startResetPasswordReq) (*startResetPasswordResp, error) {
	type startSignUpInput struct {
		Email string `validate:"required,email"`
	}

	if err := validator.GenericValidate(startSignUpInput{
		Email: req.Body.Email,
	}); err != nil {
		return nil, huma.Error400BadRequest("Invalid input", err)
	}

	user, err := s.userRepo.FindByEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, huma.Error404NotFound("user not found, consider signing up", err)
	}

	_, err = s.otp.SendOTP(model.OTPPurposePasswordReset, model.OTPChannelEmail, user.ID)
	if err != nil {
		return nil, err
	}

	resp := &startResetPasswordResp{}
	resp.Body.Message = "otp sent to email"
	return resp, nil
}

func (s *AuthService) resendResetPasswordOTP(ctx context.Context, req *resendResetPasswordOTPReq) (*resendResetPasswordOTPResp, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, huma.Error404NotFound("user not found", err)
	}

	allow, err := s.otp.ThrottleOTP(user.ID, model.OTPPurposePasswordReset)
	if err != nil {
		return nil, err
	}

	if !allow {
		return nil, huma.Error400BadRequest("can't request yet", errors.New("wait a while, can't request yet"))
	}

	_, err = s.otp.SendOTP(model.OTPPurposePasswordReset, model.OTPChannelEmail, user.ID)
	if err != nil {
		return nil, err
	}

	resp := &resendResetPasswordOTPResp{}
	resp.Body.Message = "otp sent to email again"
	return resp, nil
}

func (s *AuthService) verifyResetPasswordOTP(ctx context.Context, req *verifyResetPasswordOTPReq) (*verifyResetPasswordOTPResp, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, huma.Error404NotFound("user not found", err)
	}

	isValid, err := s.otp.VerifyOTP(model.OTPPurposePasswordReset, model.OTPChannelEmail, user.ID, req.Body.OTP, false)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, huma.Error403Forbidden("Invalid OTP", errors.New("the provided OTP is incorrect or has expired"))
	}

	resp := &verifyResetPasswordOTPResp{}
	return resp, nil
}

func (s *AuthService) resetPassword(ctx context.Context, req *resetPasswordReq) (*resetPasswordResp, error) {
	type validatePassword struct {
		Password string `validate:"required,min=8,password"`
	}

	if err := validator.GenericValidate(validatePassword{
		Password: req.Body.Password,
	}); err != nil {
		return nil, huma.Error400BadRequest("Invalid input", err)
	}

	user, err := s.userRepo.FindByEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, huma.Error404NotFound("user not found", err)
	}

	isValid, err := s.otp.VerifyOTP(model.OTPPurposePasswordReset, model.OTPChannelEmail, user.ID, req.Body.OTP, true)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, huma.Error403Forbidden("Invalid OTP", errors.New("the provided OTP is incorrect or has expired"))
	}

	err = user.SetPassword(req.Body.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError("An error occured", err)
	}

	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, huma.Error500InternalServerError("An error occured", err)
	}

	resp := &resetPasswordResp{}
	resp.Body.Message = "password reset successful, login again"
	return resp, nil
}
