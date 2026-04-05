package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	usersUseCases "github.com/felipe1496/open-wallet/internal/resources/users/usecases"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type API struct {
	googleService services.GoogleService
	usersUseCase  usersUseCases.UsersUseCases
	JWTService    services.JWTService
	authUseCases  usecases.AuthUseCases
}

func NewHandler(googleService services.GoogleService, jwtService services.JWTService, usersUseCase usersUseCases.UsersUseCases, authUseCases usecases.AuthUseCases) *API {
	return &API{
		googleService: googleService,
		usersUseCase:  usersUseCase,
		authUseCases:  authUseCases,
		JWTService:    jwtService,
	}
}

// @Summary Login with Google
// @Description Authenticates user with Google OAuth
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginGoogleRequest true "Login payload"
// @Success 200 {object} LoginGoogleResponse "User logged in"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Router /auth/login/google [post]
func (api *API) LoginGoogle(ctx *gin.Context) {
	var body LoginGoogleRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		httpErr := utils.NewHTTPError(http.StatusBadRequest, "It was not possible to process the request body")
		ctx.JSON(httpErr.StatusCode, httpErr)
		return
	}

	user, err := api.authUseCases.LoginWithGoogle(body.Code)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	access_token, err := api.JWTService.GenerateToken(user.ID)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, LoginGoogleResponse{
		Data: LoginGoogleResponseData{
			AccessToken: access_token,
			User:        user,
		},
	})
}
