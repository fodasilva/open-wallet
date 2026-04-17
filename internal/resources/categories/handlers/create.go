package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CreateOptions struct {
	Ctx      *gin.Context
	UseCases usecases.CategoriesUseCases

	UserID string
	Body   CreateCategoryRequest
}

func (o *CreateOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")

	if err := ctx.ShouldBindJSON(&o.Body); err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *CreateOptions) Validate() error {
	return nil
}

func (o *CreateOptions) Run() error {
	category, err := o.UseCases.Create(o.Ctx.Request.Context(), repository.CreateCategoryDTO{
		UserID: o.UserID,
		Name:   o.Body.Name,
		Color:  o.Body.Color,
	})

	if err != nil {
		return err
	}

	o.Ctx.JSON(http.StatusCreated, CreateCategoryResponse{
		Data: CreateCategoryResponseData{
			Category: category,
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
// @Success 201 {object} CreateCategoryResponse "Category created"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/categories [post]
func (api *API) Create(ctx *gin.Context) {
	cmd := &CreateOptions{
		UseCases: api.categoriesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
