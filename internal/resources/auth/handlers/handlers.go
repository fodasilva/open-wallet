package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	"github.com/felipe1496/open-wallet/internal/services"
)

type API struct {
	authUseCases usecases.AuthUseCases
	jwtService   services.JWTService
}

func NewHandler(authUseCases usecases.AuthUseCases, jwtService services.JWTService) *API {
	return &API{
		authUseCases: authUseCases,
		jwtService:   jwtService,
	}
}
