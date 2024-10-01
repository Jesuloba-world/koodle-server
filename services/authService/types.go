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
		Message      string      `json:"message" example:"Password set"`
		User         *model.User `json:"user" doc:"the user object"`
		AccessToken  string      `json:"accesstoken" doc:"the accesstoken for authentication"`
		RefreshToken string      `json:"refreshtoken" doc:"the refreshtoken to refresh access"`
	}
}

type loginReq struct {
	Body struct {
		Email    string `json:"email" example:"test@test.com" doc:"email of user"`
		Password string `json:"password" example:"password" doc:"password of user"`
	}
}

type loginResp struct {
	Body struct {
		Message      string      `json:"message" example:"login successful"`
		User         *model.User `json:"user" doc:"the user object"`
		AccessToken  string      `json:"accesstoken" doc:"the accesstoken for authentication"`
		RefreshToken string      `json:"refreshtoken" doc:"the refreshtoken to refresh access"`
	}
}

type refreshTokenReq struct {
	Body struct {
		RefreshToken string `json:"refreshtoken" doc:"the previous refreshtoken to refresh access"`
	}
}

type refreshTokenResp struct {
	Body struct {
		Message      string `json:"message" example:"refresh successful"`
		AccessToken  string `json:"accesstoken" doc:"the accesstoken for authentication"`
		RefreshToken string `json:"refreshtoken" doc:"the refreshtoken to refresh access"`
	}
}

type startResetPasswordReq struct {
	Body struct {
		Email string `json:"email" example:"test@test.com" doc:"email of user, otp will be sent to this email"`
	}
}
type startResetPasswordResp struct {
	Body struct {
		Message string `json:"message" example:"continue"`
	}
}

type resendResetPasswordOTPReq struct {
	Body struct {
		Email string `json:"email" example:"test@test.com" doc:"email of user, otp will be sent to this email"`
	}
}

type resendResetPasswordOTPResp struct {
	Body struct {
		Message string `json:"message" example:"OTP sent to email"`
	}
}

type verifyResetPasswordOTPReq struct {
	Body struct {
		Email string `json:"email" example:"test@test.com" doc:"email of user"`
		OTP   string `json:"otp" example:"123456" doc:"otp sent to email"`
	}
}

type verifyResetPasswordOTPResp struct {
	Body struct {
		Message string `json:"message" example:"OTP verified"`
	}
}

type resetPasswordReq struct {
	Body struct {
		Email    string `json:"email" example:"test@test.com" doc:"email of user"`
		Password string `json:"password" example:"password" doc:"new password of user"`
		OTP      string `json:"otp" example:"123456" doc:"otp sent to email"`
	}
}

type resetPasswordResp struct {
	Body struct {
		Message string `json:"message" example:"Password reset successful"`
	}
}
