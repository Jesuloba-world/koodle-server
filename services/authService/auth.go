package authservice

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/exp/slog"

	humagroup "github.com/Jesuloba-world/koodle-server/lib/humaGroup"

)

type AuthService struct {
	api *humagroup.HumaGroup
}

func NewAuthService(api huma.API) *AuthService {
	return &AuthService{
		api: humagroup.NewHumaGroup(api, "/auth", []string{"Authentication"}),
	}
}

func (s *AuthService) RegisterRoutes() {
	humagroup.Post(s.api, "/startsignup", s.StartSignUp, "Start signup process")
}

func (s *AuthService) StartSignUp(ctx context.Context, req *startSignUpReq) (*startSignUpResp, error) {
	slog.Info(req.Body.Email)

	resp := &startSignUpResp{}
	resp.Body.Success = true
	return resp, nil
}
