package handlers

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	_ "github.com/felipe1496/open-wallet/internal/util/httputil"
)

type DeleteOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.CategoriesUseCases

	ID     string
	UserID string
}

func (o *DeleteOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.ID = r.PathValue("category_id")
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)

	return nil
}

func (o *DeleteOptions) Validate() error {
	return nil
}

func (o *DeleteOptions) Run() error {
	err := o.UseCases.DeleteByID(o.R.Context(), o.ID, o.UserID)
	if err != nil {
		return err
	}

	o.W.WriteHeader(http.StatusNoContent)
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
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 404 {object} httputil.HTTPError "Not found"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Failure 503 {string} string "Service Unavailable"
// @Router /api/v1/categories/{category_id} [delete]
func (api *API) DeleteByID(w http.ResponseWriter, r *http.Request) {
	cmd := &DeleteOptions{
		UseCases: api.categoriesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
