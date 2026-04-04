package auth

import (
	
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/users"
	
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/gin-gonic/gin"
)

type API struct {
	googleService services.GoogleService
	usersUseCase  users.UsersUseCase
	JWTService    services.JWTService
	authUseCase   AuthUseCase
}

func NewHandler(googleService services.GoogleService, jwtService services.JWTService, usersUseCase users.UsersUseCase, authUseCase AuthUseCase) *API {
	return &API{
		googleService: googleService,
		usersUseCase:  usersUseCase,
		authUseCase:   authUseCase,
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

	user, err := api.authUseCase.LoginWithGoogle(body.Code)

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
