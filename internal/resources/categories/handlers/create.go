package handlers

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type CreateOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.CategoriesUseCases

	UserID string
	Body   CreateCategoryRequest
}

func (o *CreateOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)

	if err := httputil.BindJSON(r, &o.Body); err != nil {
		return httputil.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *CreateOptions) Validate() error {
	if len(o.Body.Name) == 0 {
		return httputil.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	if o.Body.Color == "" {
		return httputil.NewHTTPError(http.StatusBadRequest, "color is required")
	}
	return nil
}

func (o *CreateOptions) Run() error {
	category, err := o.UseCases.Create(o.R.Context(), repository.CreateCategoryDTO{
		UserID: o.UserID,
		Name:   o.Body.Name,
		Color:  o.Body.Color,
	})

	if err != nil {
		return err
	}

	httputil.JSON(o.W, http.StatusCreated, util.ResponseData[CreateCategoryResponseData]{
		Data: CreateCategoryResponseData{
			Category: MapCategoryResource(category),
		},
	})
	return nil
}

// @Summary Create a category
// @Description Create a category
// @ID v1CreateCategory
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateCategoryRequest true "Category payload"
// @Success 201 {object} util.ResponseData[CreateCategoryResponseData] "Category created"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Failure 503 {string} string "Service Unavailable"
// @Router /api/v1/categories [post]
func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	cmd := &CreateOptions{
		UseCases: api.categoriesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
