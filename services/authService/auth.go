package authservice

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bun"

	humagroup "github.com/Jesuloba-world/koodle-server/lib/humaGroup"
	"github.com/Jesuloba-world/koodle-server/lib/validator"
	"github.com/Jesuloba-world/koodle-server/model"
	userrepo "github.com/Jesuloba-world/koodle-server/repo/user"
	"github.com/Jesuloba-world/koodle-server/services/otpService"

)

type AuthService struct {
	api      *humagroup.HumaGroup
	userRepo *userrepo.UserRepo
	otp      *otpService.OTPService
}

func NewAuthService(api huma.API, db *bun.DB, otpService *otpService.OTPService, userRepo *userrepo.UserRepo) *AuthService {
	return &AuthService{
		api:      humagroup.NewHumaGroup(api, "/auth", []string{"Authentication"}),
		userRepo: userRepo,
		otp:      otpService,
	}
}

func (s *AuthService) RegisterRoutes() {
	humagroup.Post(s.api, "/startsignup", s.StartSignUp, "Start signup process")
}

func (s *AuthService) StartSignUp(ctx context.Context, req *startSignUpReq) (*startSignUpResp, error) {
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
			return nil, huma.Error404NotFound("not found", err)
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
	resp.Body.Success = true
	return resp, nil
}
