package authservice

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Jesuloba-world/koodle-server/lib/validator"
	"github.com/Jesuloba-world/koodle-server/model"
)

func (s *AuthService) startSignUp(ctx context.Context, req *startSignUpReq) (*startSignUpResp, error) {
	type startSignUpInput struct {
		Email string `validate:"required,email"`
	}

	if err := validator.GenericValidate(startSignUpInput{
		Email: req.Body.Email,
	}); err != nil {
		return nil, huma.Error400BadRequest("Invalid input", err)
	}

	var user model.User

	exists, passwordSet, err := s.userRepo.CheckUserExists(ctx, req.Body.Email)
	if err != nil {
		return nil, huma.Error403Forbidden("an error occured", err)
	}
	if exists {
		if passwordSet {
			return nil, huma.Error403Forbidden("user already exists", errors.New("user already exists, consider logging in instead"))
		}
		oldUser, err := s.userRepo.FindByEmail(ctx, req.Body.Email)
		if err != nil {
			return nil, huma.Error404NotFound("user not found", err)
		}
		user = *oldUser
	} else {
		newUser := model.User{
			Email:         req.Body.Email,
			EmailVerified: false,
		}

		err = s.userRepo.CreateUser(ctx, &newUser)
		if err != nil {
			return nil, huma.Error500InternalServerError("an error occured", err)
		}

		user = newUser
	}

	_, err = s.otp.SendOTP(model.OTPPurposeEmailVerification, model.OTPChannelEmail, user.ID)
	if err != nil {
		return nil, err
	}

	resp := &startSignUpResp{}
	resp.Body.Message = "Signup started"
	return resp, nil
}

func (s *AuthService) resendEmailVerificationOTP(ctx context.Context, req *resendEmailVerificationOTPReq) (*resendEmailVerificationOTPResp, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, huma.Error404NotFound("user not found", err)
	}

	allow, err := s.otp.ThrottleOTP(user.ID, model.OTPPurposeEmailVerification)
	if err != nil {
		return nil, err
	}

	if !allow {
		return nil, huma.Error400BadRequest("can't request yet", errors.New("wait a while, can't request yet"))
	}

	_, err = s.otp.SendOTP(model.OTPPurposeEmailVerification, model.OTPChannelEmail, user.ID)
	if err != nil {
		return nil, err
	}

	resp := &resendEmailVerificationOTPResp{}
	resp.Body.Message = "OTP sent to email"
	return resp, nil
}

func (s *AuthService) verifyEmail(ctx context.Context, req *verifyEmailReq) (*verifyEmailResp, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, huma.Error404NotFound("user not found", err)
	}

	isValid, err := s.otp.VerifyOTP(model.OTPPurposeEmailVerification, model.OTPChannelEmail, user.ID, req.Body.OTP, false)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, huma.Error403Forbidden("Invalid OTP", errors.New("the provided OTP is incorrect or has expired"))
	}

	user.EmailVerified = true

	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	resp := &verifyEmailResp{}
	resp.Body.Message = "Email verified successfully"
	return resp, nil
}

func (s *AuthService) setPassword(ctx context.Context, req *setPasswordReq) (*setPasswordResp, error) {
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

	isValid, err := s.otp.VerifyOTP(model.OTPPurposeEmailVerification, model.OTPChannelEmail, user.ID, req.Body.OTP, false)
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

	access, refresh, err := s.token.GenerateTokens(user.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("could not generate tokens", err)
	}

	resp := &setPasswordResp{}
	resp.Body.Message = "User signup completed"
	resp.Body.User = *user
	resp.Body.AccessToken = access
	resp.Body.RefreshToken = refresh
	return resp, nil
}
