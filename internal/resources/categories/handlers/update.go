package handlers

import (
	"net/http"
	"slices"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type UpdateOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.CategoriesUseCases

	ID         string
	UserID     string
	PassedKeys []string
	Body       UpdateCategoryRequest
	Payload    repository.UpdateCategoryDTO
}

func (o *UpdateOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.ID = r.PathValue("category_id")
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)

	passedKeys, err := httputil.GetJSONKeys(r)
	if err != nil {
		return httputil.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	o.PassedKeys = passedKeys

	if len(o.PassedKeys) == 0 {
		return httputil.NewHTTPError(http.StatusBadRequest, "At least one field must be provided for update")
	}

	if err := httputil.BindJSON(r, &o.Body); err != nil {
		return httputil.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *UpdateOptions) Validate() error {
	if slices.Contains(o.PassedKeys, "name") {
		if o.Body.Name == nil {
			return httputil.NewHTTPError(http.StatusBadRequest, "name cannot be null")
		}
		if len(*o.Body.Name) == 0 {
			return httputil.NewHTTPError(http.StatusBadRequest, "name cannot be empty")
		}
		o.Payload.Name = util.NewValue(*o.Body.Name)
	}

	if slices.Contains(o.PassedKeys, "color") {
		if o.Body.Color == nil {
			return httputil.NewHTTPError(http.StatusBadRequest, "color cannot be null")
		}
		if len(*o.Body.Color) == 0 {
			return httputil.NewHTTPError(http.StatusBadRequest, "color cannot be empty")
		}
		o.Payload.Color = util.NewValue(*o.Body.Color)
	}

	return nil
}

func (o *UpdateOptions) Run() error {
	category, err := o.UseCases.Update(o.R.Context(), o.ID, o.UserID, o.Payload)
	if err != nil {
		return err
	}

	httputil.JSON(o.W, http.StatusOK, util.ResponseData[UpdateCategoryResponseData]{
		Data: UpdateCategoryResponseData{
			Category: MapCategoryResource(category),
		},
	})
	return nil
}

// @Summary Update Category By ID
// @Description Update a category
// @ID v1UpdateCategory
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category_id path string true "category ID"
// @Param body body UpdateCategoryRequest true "Category payload"
// @Success 200 {object} util.ResponseData[UpdateCategoryResponseData] "Category updated"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 404 {object} httputil.HTTPError "Not found"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Router /api/v1/categories/{category_id} [patch]
func (api *API) Update(w http.ResponseWriter, r *http.Request) {
	cmd := &UpdateOptions{
		UseCases: api.categoriesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
