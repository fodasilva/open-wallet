package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	usersUseCases "github.com/felipe1496/open-wallet/internal/resources/users/usecases"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"
)

var (
	GoogleAuthFailedErr       = utils.NewHTTPError(http.StatusUnauthorized, "authentication with Google failed")
	GoogleDintProvideEmailErr = utils.NewHTTPError(http.StatusUnauthorized, "google did not provide an email")
	JwtGenErr                 = utils.NewHTTPError(http.StatusUnauthorized, "failed to generate JWT token")
	GoogleEmailNotVerifiedErr = utils.NewHTTPError(http.StatusUnauthorized, "Google email not verified")
)

type AuthUseCases interface {
	LoginWithGoogle(code string) (repository.User, error)
}

type AuthUseCasesImpl struct {
	googleService services.GoogleService
	usersUseCase  usersUseCases.UsersUseCases
}

func NewAuthUseCases(googleService services.GoogleService, usersUseCase usersUseCases.UsersUseCases) AuthUseCases {
	return &AuthUseCasesImpl{
		googleService: googleService,
		usersUseCase:  usersUseCase,
	}
}
