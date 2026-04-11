package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CreateLoginWithGoogleOptions struct {
	Ctx        *gin.Context
	UseCases   usecases.AuthUseCases
	JWTService services.JWTService

	Body LoginGoogleRequest
}

func (o *CreateLoginWithGoogleOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx

	if err := ctx.ShouldBindJSON(&o.Body); err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, "it was not possible to process the request body")
	}

	return nil
}

func (o *CreateLoginWithGoogleOptions) Validate() error {
	return nil
}

func (o *CreateLoginWithGoogleOptions) Run() error {
	user, err := o.UseCases.LoginWithGoogle(o.Body.Code)
	if err != nil {
		return err
	}

	accessToken, err := o.JWTService.GenerateToken(user.ID)
	if err != nil {
		return err
	}

	o.Ctx.JSON(http.StatusOK, LoginGoogleResponse{
		Data: LoginGoogleResponseData{
			AccessToken: accessToken,
			User:        user,
		},
	})

	return nil
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
// @Router /api/v1/auth/login/google [post]
func (api *API) CreateLoginWithGoogle(ctx *gin.Context) {
	cmd := &CreateLoginWithGoogleOptions{
		UseCases:   api.authUseCases,
		JWTService: api.jwtService,
	}
	utils.RunCommand(ctx, cmd)
}
