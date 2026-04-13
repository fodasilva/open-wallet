package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

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
	reqCtx := querybuilder.WithBuilder(o.Ctx, o.Builder)
	categoriesList, err := o.UseCases.List(reqCtx)
	if err != nil {
		return err
	}

	countCtx := querybuilder.WithBuilder(o.Ctx, querybuilder.ForCount(o.Builder))
	count, err := o.UseCases.Count(countCtx)
	if err != nil {
		return err
	}

	nextPage := len(categoriesList) > o.PerPage
	totalPages := (count + o.PerPage - 1) / o.PerPage

	if nextPage {
		categoriesList = categoriesList[:len(categoriesList)-1]
	}

	o.Ctx.JSON(http.StatusOK, ListCategoriesResponse{
		Data: ListCategoriesResponseData{
			Categories: categoriesList,
		},
		Query: utils.QueryMeta{
			NextPage:   nextPage,
			Page:       o.Page,
			PerPage:    o.PerPage,
			TotalItems: count,
			TotalPages: totalPages,
		},
	})

	return nil
}

// @Summary List categories
// @Description List categories
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param order_by query string false "Sort field" example(name:asc,created_at:desc)
// @Param filter query string false "Category filter"
// @Param name query string false "A category name to filter by"
// @Success 200 {object} ListCategoriesResponse "List of categories"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/categories [get]
func (api *API) List(ctx *gin.Context) {
	cmd := &ListOptions{
		UseCases: api.categoriesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
