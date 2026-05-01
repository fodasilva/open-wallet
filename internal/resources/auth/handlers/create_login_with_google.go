package handlers

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type CreateLoginWithGoogleOptions struct {
	W          http.ResponseWriter
	R          *http.Request
	UseCases   usecases.AuthUseCases
	JWTService services.JWTService

	Body LoginGoogleRequest
}

func (o *CreateLoginWithGoogleOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r

	if err := httputil.BindJSON(r, &o.Body); err != nil {
		return util.NewHTTPError(http.StatusBadRequest, "it was not possible to process the request body")
	}

	return nil
}

func (o *CreateLoginWithGoogleOptions) Validate() error {
	if len(o.Body.Code) == 0 {
		return util.NewHTTPError(http.StatusBadRequest, "code is required")
	}
	return nil
}

func (o *CreateLoginWithGoogleOptions) Run() error {
	user, err := o.UseCases.LoginWithGoogle(o.R.Context(), o.Body.Code)
	if err != nil {
		return err
	}

	accessToken, err := o.JWTService.GenerateToken(user.ID)
	if err != nil {
		return err
	}

	httputil.JSON(o.W, http.StatusOK, util.ResponseData[LoginGoogleResponseData]{
		Data: LoginGoogleResponseData{
			AccessToken: accessToken,
			User:        MapUserResource(user),
		},
	})

	return nil
}

// @Summary Login with Google
// @Description Authenticates user with Google OAuth
// @ID v1LoginWithGoogle
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginGoogleRequest true "Login payload"
// @Success 200 {object} util.ResponseData[LoginGoogleResponseData] "User logged in"
// @Failure 400 {object} util.HTTPError "Bad request"
// @Failure 401 {object} util.HTTPError "Unauthorized"
// @Router /api/v1/auth/login/google [post]
func (api *API) CreateLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	cmd := &CreateLoginWithGoogleOptions{
		UseCases:   api.authUseCases,
		JWTService: api.jwtService,
	}
	util.RunCommand(w, r, cmd)
}
