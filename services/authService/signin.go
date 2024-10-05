package authservice

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
)

func (s *AuthService) login(ctx context.Context, req *loginReq) (*loginResp, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, huma.Error404NotFound("Invalid email or password", err)
	}

	if !user.CheckPassword(req.Body.Password) {
		return nil, huma.Error401Unauthorized("Invalid email or password", errors.New("invalid email or password"))
	}

	access, refresh, err := s.token.GenerateTokens(user.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("could not generate tokens", err)
	}

	resp := &loginResp{}
	resp.Body.Message = "Login successful"
	resp.Body.User = user
	resp.Body.AccessToken = access
	resp.Body.RefreshToken = refresh
	return resp, nil
}

func (s *AuthService) refreshToken(ctx context.Context, req *refreshTokenReq) (*refreshTokenResp, error) {
	access, refresh, err := s.token.RefreshToken(req.Body.RefreshToken)
	if err != nil {
		return nil, huma.Error500InternalServerError("could not generate tokens", err)
	}

	resp := &refreshTokenResp{}
	resp.Body.Message = "Token refreshed successfully"
	resp.Body.AccessToken = access
	resp.Body.RefreshToken = refresh
	return resp, nil
}
