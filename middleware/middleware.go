package custommiddleware

import (
	"github.com/danielgtaylor/huma/v2"

	tokenservice "github.com/Jesuloba-world/koodle-server/services/tokenService"
)

type Middleware struct {
	api          huma.API
	tokenService *tokenservice.TokenService
}

func MakeMiddleware(api huma.API, tokenservice *tokenservice.TokenService) *Middleware {
	return &Middleware{
		api:          api,
		tokenService: tokenservice,
	}
}
