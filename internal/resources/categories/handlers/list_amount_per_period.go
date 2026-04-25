package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
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
	Ctx      *gin.Context
	UseCases usecases.CategoriesUseCases
	UserID   string
	Period   string
	Builder  *querybuilder.Builder
	Page     int
	PerPage  int
}

func (o *ListAmountPerPeriodOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")
	o.Period = ctx.Param("period")
	o.Builder = ctx.MustGet("query_builder").(*querybuilder.Builder).
		And("user_id", "eq", o.UserID)
	o.Page = o.Ctx.GetInt("page")
	o.PerPage = o.Ctx.GetInt("per_page")

	return nil
}

func (o *ListAmountPerPeriodOptions) Validate() error {
	return nil
}

func (o *ListAmountPerPeriodOptions) Run() error {
	reqCtx := querybuilder.WithBuilder(o.Ctx.Request.Context(), o.Builder)
	categories, err := o.UseCases.ListCategoryAmountPerPeriod(reqCtx, o.Period)
	if err != nil {
		return err
	}

	countCtx := querybuilder.WithBuilder(o.Ctx.Request.Context(), querybuilder.ForCount(o.Builder))
	count, err := o.UseCases.CountCategoryAmountPerPeriod(countCtx, o.Period)
	if err != nil {
		return err
	}

	categoriesResource := make([]CategoryAmountPerPeriodResource, len(categories))
	for i, c := range categories {
		categoriesResource[i] = MapCategoryAmountPerPeriodResource(c)
	}

	o.Ctx.JSON(http.StatusOK, utils.PaginatedResponse[ListCategoryAmountPerPeriodResponseData]{
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
// @Success 200 {object} utils.PaginatedResponse[ListCategoryAmountPerPeriodResponseData] "List of categories with amount per period"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/categories/{period} [get]
func (api *API) ListCategoryAmountPerPeriod(ctx *gin.Context) {
	cmd := &ListAmountPerPeriodOptions{
		UseCases: api.categoriesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
