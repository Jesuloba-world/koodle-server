package authservice

import "github.com/Jesuloba-world/koodle-server/model"

type startSignUpReq struct {
	Body struct {
		Email string `json:"email" doc:"email of new user, otp will be sent to this email for verification" example:"test@test.com"`
	}
}

type startSignUpResp struct {
	Body struct {
		Message string `json:"message" example:"Signup started"`
	}
}

type resendEmailVerificationOTPReq struct {
	Body struct {
		Email string `json:"email" example:"test@test.com" doc:"resend otp to email"`
	}
}

type resendEmailVerificationOTPResp struct {
	Body struct {
		Message string `json:"message" example:"OTP sent to email"`
	}
}

type verifyEmailReq struct {
	Body struct {
		Email string `json:"email" example:"test@test.com" doc:"email of user"`
		OTP   string `json:"otp" example:"123456" doc:"otp sent to email"`
	}
}

type verifyEmailResp struct {
	Body struct {
		Message string `json:"message" example:"Email verified"`
	}
}

type setPasswordReq struct {
	Body struct {
		Email    string `json:"email" example:"test@test.com" doc:"email of user"`
		Password string `json:"password" example:"password" doc:"new password of user"`
		OTP      string `json:"otp" example:"123456" doc:"otp sent to email"`
	}
}

type setPasswordResp struct {
	Body struct {
		Message      string     `json:"message" example:"Password set"`
		User         model.User `json:"user" doc:"the user object"`
		AccessToken  string     `json:"accesstoken" doc:"the accesstoken for authentication"`
		RefreshToken string     `json:"refreshtoken" doc:"the refreshtoken to refresh access"`
	}
}
