package authservice

type startSignUpReq struct {
	Body struct {
		Email string `json:"email" doc:"email of new user, otp wiill be sent to this email for verification" example:"test@test.com"`
	}
}

type startSignUpResp struct {
	Body struct {
		Success bool `json:"success" doc:"true if sign up started" example:"true"`
	}
}
