package custommiddleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	tokenservice "github.com/Jesuloba-world/koodle-server/services/tokenService"
)

var (
	ErrNoAuthToken  = errors.New("no auth token provided")
	ErrInvalidToken = errors.New("invalid auth token")
	ErrInvalidID    = errors.New("invalid id")
)

func (m *Middleware) Auth(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized,
			"Unauthenticated request", ErrNoAuthToken,
		)
		return
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	claims, tokenType, err := m.tokenService.ValidateToken(accessToken)
	if err != nil || tokenType != tokenservice.AccessTokenType {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized,
			"Invalid token", ErrInvalidToken,
		)
		return
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		huma.WriteErr(m.api, ctx, http.StatusUnauthorized,
			"Invalid id", ErrInvalidID,
		)
		return
	}

	ctx = huma.WithValue(ctx, "user_id", userId)

	next(ctx)
}
