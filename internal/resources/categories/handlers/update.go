package handlers

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type UpdateOptions struct {
	Ctx      *gin.Context
	UseCases usecases.CategoriesUseCases

	ID         string
	UserID     string
	PassedKeys []string
	Body       UpdateCategoryRequest
	Payload    repository.UpdateCategoryDTO
}

func (o *UpdateOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.ID = ctx.Param("category_id")
	o.UserID = ctx.GetString("user_id")

	passedKeys, err := utils.GetJSONKeys(ctx)
	if err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	o.PassedKeys = passedKeys

	if len(o.PassedKeys) == 0 {
		return utils.NewHTTPError(http.StatusBadRequest, "At least one field must be provided for update")
	}

	if err := ctx.ShouldBindJSON(&o.Body); err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *UpdateOptions) Validate() error {
	if slices.Contains(o.PassedKeys, "name") {
		if o.Body.Name == nil {
			return utils.NewHTTPError(http.StatusBadRequest, "name cannot be null")
		}
		o.Payload.Name = utils.NewValue(*o.Body.Name)
	}

	if slices.Contains(o.PassedKeys, "color") {
		if o.Body.Color == nil {
			return utils.NewHTTPError(http.StatusBadRequest, "color cannot be null")
		}
		o.Payload.Color = utils.NewValue(*o.Body.Color)
	}

	return nil
}

func (o *UpdateOptions) Run() error {
	category, err := o.UseCases.Update(o.ID, o.UserID, o.Payload)
	if err != nil {
		return err
	}

	o.Ctx.JSON(http.StatusOK, UpdateCategoryResponse{
		Data: UpdateCategoryResponseData{
			Category: category,
		},
	})
	return nil
}

// @Summary Update Category By ID
// @Description Update a category
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category_id path string true "category ID"
// @Param body body UpdateCategoryRequest true "Category payload"
// @Success 200 {object} UpdateCategoryResponse "Category updated"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 404 {object} utils.HTTPError "Not found"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/categories/{category_id} [patch]
func (api *API) Update(ctx *gin.Context) {
	cmd := &UpdateOptions{
		UseCases: api.categoriesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
