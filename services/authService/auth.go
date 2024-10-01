package authservice

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bun"

	humagroup "github.com/Jesuloba-world/koodle-server/lib/humaGroup"
	userrepo "github.com/Jesuloba-world/koodle-server/repo/user"
	"github.com/Jesuloba-world/koodle-server/services/otpService"
	tokenservice "github.com/Jesuloba-world/koodle-server/services/tokenService"
)

type AuthService struct {
	api      *humagroup.HumaGroup
	userRepo *userrepo.UserRepo
	otp      *otpService.OTPService
	token    *tokenservice.TokenService
}

func NewAuthService(api huma.API, db *bun.DB, otpService *otpService.OTPService, userRepo *userrepo.UserRepo, tokenService *tokenservice.TokenService) *AuthService {
	return &AuthService{
		api:      humagroup.NewHumaGroup(api, "/auth", []string{"Authentication"}),
		userRepo: userRepo,
		otp:      otpService,
		token:    tokenService,
	}
}

func (s *AuthService) RegisterRoutes() {
	humagroup.Post(s.api, "/startsignup", s.startSignUp, "Start signup process")
	humagroup.Post(s.api, "/resendemailverificationotp", s.resendEmailVerificationOTP, "Resend email verification OTP")
	humagroup.Post(s.api, "/verifyemail", s.verifyEmail, "Verify email address with OTP")
	humagroup.Post(s.api, "/setPassword", s.setPassword, "Set Password for user")
	humagroup.Post(s.api, "/login", s.login, "Login user")
	humagroup.Post(s.api, "/refreshtoken", s.refreshToken, "Refresh token")
	humagroup.Post(s.api, "/startresetpassword", s.startResetPassword, "Start Reset Password")
	humagroup.Post(s.api, "/resendresetpasswordotp", s.resendResetPasswordOTP, "Resend reset password OTP")
	humagroup.Post(s.api, "/verifyresetpasswordotp", s.verifyResetPasswordOTP, "Verify reset password OTP")
	humagroup.Post(s.api, "/resetpassword", s.resetPassword, "Reset password")
}
