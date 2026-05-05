package handlers

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

// @gen_swagger_filter
var PeriodCategoriesFilterConfig = querybuilder.ParseConfig{
	AllowedFields: map[string]querybuilder.FieldConfig{
		"name":         {AllowedOperators: []string{"eq", "like", "in"}},
		"color":        {AllowedOperators: []string{"eq", "in"}},
		"total_amount": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"period":       {AllowedOperators: []string{"eq", "in"}},
		"id":           {AllowedOperators: []string{"eq", "in"}},
		"user_id":      {AllowedOperators: []string{"eq", "in"}},
	},
	AllowedSortFields: []string{"name", "total_amount", "period", "id"},
}

type ListAmountPerPeriodOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.CategoriesUseCases
	UserID   string
	Period   string
	Builder  *querybuilder.Builder
	Page     int
	PerPage  int
}

func (o *ListAmountPerPeriodOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)
	o.Period = r.PathValue("period")
	o.Builder = querybuilder.Get(r.Context()).
		And("user_id", "eq", o.UserID)
	o.Page = util.GetInt(r.Context(), util.ContextKeyPage)
	o.PerPage = util.GetInt(r.Context(), util.ContextKeyPerPage)

	return nil
}

func (o *ListAmountPerPeriodOptions) Validate() error {
	return nil
}

func (o *ListAmountPerPeriodOptions) Run() error {
	reqCtx := querybuilder.WithBuilder(o.R.Context(), o.Builder)
	categories, err := o.UseCases.ListCategoryAmountPerPeriod(reqCtx, o.Period)
	if err != nil {
		return err
	}

	countCtx := querybuilder.WithBuilder(o.R.Context(), querybuilder.ForCount(o.Builder))
	count, err := o.UseCases.CountCategoryAmountPerPeriod(countCtx, o.Period)
	if err != nil {
		return err
	}

	categoriesResource := make([]CategoryAmountPerPeriodResource, len(categories))
	for i, c := range categories {
		categoriesResource[i] = MapCategoryAmountPerPeriodResource(c)
	}

	httputil.JSON(o.W, http.StatusOK, util.PaginatedResponse[ListCategoryAmountPerPeriodResponseData]{
		Data: ListCategoryAmountPerPeriodResponseData{
			Categories: categoriesResource,
		},
		Query: querybuilder.BuildMetadata(o.Page, o.PerPage, count),
	})
	return nil
}

// @Summary List categories with amount per period
// @Description List categories with amount per period
// @ID v1ListCategoryAmountPerPeriod
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period path string true "period"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param filter query string false "Filter expression. \n- Allowed fields & ops:\n  - color: eq, in\n  - id: eq, in\n  - name: eq, like, in\n  - period: eq, in\n  - total_amount: eq, gt, gte, lt, lte\n  - user_id: eq, in\n"
// @Param order_by query string false "Sort field. \n- Allowed: name, total_amount, period, id" example(name:asc)
// @Success 200 {object} util.PaginatedResponse[ListCategoryAmountPerPeriodResponseData] "List of categories with amount per period"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Failure 503 {string} string "Service Unavailable"
// @Router /api/v1/categories/{period} [get]
func (api *API) ListCategoryAmountPerPeriod(w http.ResponseWriter, r *http.Request) {
	cmd := &ListAmountPerPeriodOptions{
		UseCases: api.categoriesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
