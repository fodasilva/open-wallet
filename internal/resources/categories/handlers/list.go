package handlers

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

// @gen_swagger_filter
var CategoriesFilterConfig = querybuilder.ParseConfig{
	AllowedFields: map[string]querybuilder.FieldConfig{
		"name":       {AllowedOperators: []string{"eq", "like", "in"}},
		"color":      {AllowedOperators: []string{"eq", "in"}},
		"created_at": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"id":         {AllowedOperators: []string{"eq", "in"}},
		"user_id":    {AllowedOperators: []string{"eq", "in"}},
	},
	AllowedSortFields: []string{"name", "created_at", "id"},
}

type ListOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.CategoriesUseCases
	UserID   string
	Builder  *querybuilder.Builder
	Page     int
	PerPage  int
}

func (o *ListOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)
	nameFilter := r.URL.Query().Get("name")
	o.Builder = querybuilder.Get(r.Context()).
		And("user_id", "eq", o.UserID)
	o.Page = util.GetInt(r.Context(), util.ContextKeyPage)
	o.PerPage = util.GetInt(r.Context(), util.ContextKeyPerPage)

	if nameFilter != "" {
		o.Builder.And("name", "like", nameFilter)
	}

	return nil
}

func (o *ListOptions) Validate() error {
	return nil
}

func (o *ListOptions) Run() error {
	reqCtx := querybuilder.WithBuilder(o.R.Context(), o.Builder)
	categoriesList, err := o.UseCases.List(reqCtx)
	if err != nil {
		return err
	}

	countCtx := querybuilder.WithBuilder(o.R.Context(), querybuilder.ForCount(o.Builder))
	count, err := o.UseCases.Count(countCtx)
	if err != nil {
		return err
	}

	categoriesResource := make([]CategoryResource, len(categoriesList))
	for i, c := range categoriesList {
		categoriesResource[i] = MapCategoryResource(c)
	}

	httputil.JSON(o.W, http.StatusOK, util.PaginatedResponse[ListCategoriesResponseData]{
		Data: ListCategoriesResponseData{
			Categories: categoriesResource,
		},
		Query: querybuilder.BuildMetadata(o.Page, o.PerPage, count),
	})

	return nil
}

// @Summary List categories
// @Description List categories
// @ID v1ListCategories
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param order_by query string false "Sort field. \n- Allowed: name, created_at, id" example(name:asc)
// @Param filter query string false "Filter expression. \n- Allowed fields & ops:\n  - color: eq, in\n  - created_at: eq, gt, gte, lt, lte\n  - id: eq, in\n  - name: eq, like, in\n  - user_id: eq, in\n"
// @Param name query string false "A category name to filter by"
// @Success 200 {object} util.PaginatedResponse[ListCategoriesResponseData] "List of categories"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Router /api/v1/categories [get]
func (api *API) List(w http.ResponseWriter, r *http.Request) {
	cmd := &ListOptions{
		UseCases: api.categoriesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
