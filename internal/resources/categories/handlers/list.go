package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
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
	Ctx      *gin.Context
	UseCases usecases.CategoriesUseCases
	UserID   string
	Builder  *querybuilder.Builder
	Page     int
	PerPage  int
}

func (o *ListOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")
	nameFilter := ctx.Query("name")
	o.Builder = ctx.MustGet("query_builder").(*querybuilder.Builder).
		And("user_id", "eq", o.UserID)
	o.Page = o.Ctx.GetInt("page")
	o.PerPage = o.Ctx.GetInt("per_page")

	if nameFilter != "" {
		o.Builder.And("name", "like", nameFilter)
	}

	return nil
}

func (o *ListOptions) Validate() error {
	return nil
}

func (o *ListOptions) Run() error {
	reqCtx := querybuilder.WithBuilder(o.Ctx.Request.Context(), o.Builder)
	categoriesList, err := o.UseCases.List(reqCtx)
	if err != nil {
		return err
	}

	countCtx := querybuilder.WithBuilder(o.Ctx.Request.Context(), querybuilder.ForCount(o.Builder))
	count, err := o.UseCases.Count(countCtx)
	if err != nil {
		return err
	}

	categoriesResource := make([]CategoryResource, len(categoriesList))
	for i, c := range categoriesList {
		categoriesResource[i] = MapCategoryResource(c)
	}

	o.Ctx.JSON(http.StatusOK, utils.PaginatedResponse[ListCategoriesResponseData]{
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
// @Success 200 {object} utils.PaginatedResponse[ListCategoriesResponseData] "List of categories"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/categories [get]
func (api *API) List(ctx *gin.Context) {
	cmd := &ListOptions{
		UseCases: api.categoriesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
