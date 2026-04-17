package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type DeleteOptions struct {
	Ctx      *gin.Context
	UseCases usecases.CategoriesUseCases

	ID     string
	UserID string
}

func (o *DeleteOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.ID = ctx.Param("category_id")
	o.UserID = ctx.GetString("user_id")

	return nil
}

func (o *DeleteOptions) Validate() error {
	return nil
}

func (o *DeleteOptions) Run() error {
	err := o.UseCases.DeleteByID(o.Ctx.Request.Context(), o.ID, o.UserID)
	if err != nil {
		return err
	}

	o.Ctx.Status(http.StatusNoContent)
	return nil
}

// @Summary Delete Category By ID
// @Description Delete a category
// @ID v1DeleteCategory
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category_id path string true "category ID"
// @Success 204 "Category deleted"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 404 {object} utils.HTTPError "Not found"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/categories/{category_id} [delete]
func (api *API) DeleteByID(ctx *gin.Context) {
	cmd := &DeleteOptions{
		UseCases: api.categoriesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
