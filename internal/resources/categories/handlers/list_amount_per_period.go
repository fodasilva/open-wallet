package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type ListAmountPerPeriodOptions struct {
	Ctx       *gin.Context
	UseCases  usecases.CategoriesUseCases
	UserID    string
	Period    string
	QueryOpts *utils.QueryOptsBuilder
	Page      int
	PerPage   int
}

func (o *ListAmountPerPeriodOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")
	o.Period = ctx.Param("period")
	o.QueryOpts = ctx.MustGet("query_opts").(*utils.QueryOptsBuilder).
		And("user_id", "eq", o.UserID)
	o.Page = o.Ctx.GetInt("page")
	o.PerPage = o.Ctx.GetInt("per_page")

	return nil
}

func (o *ListAmountPerPeriodOptions) Validate() error {
	return nil
}

func (o *ListAmountPerPeriodOptions) Run() error {
	categories, err := o.UseCases.ListCategoryAmountPerPeriod(o.Period, o.QueryOpts)
	if err != nil {
		return err
	}

	count, err := o.UseCases.CountCategoryAmountPerPeriod(o.Period, utils.ForCount(o.QueryOpts))
	if err != nil {
		return err
	}

	nextPage := len(categories) > o.PerPage
	totalPages := (count + o.PerPage - 1) / o.PerPage

	if nextPage {
		categories = categories[:len(categories)-1]
	}

	o.Ctx.JSON(http.StatusOK, ListCategoryAmountPerPeriodResponse{
		Data: ListCategoryAmountPerPeriodResponseData{
			Categories: categories,
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

// @Summary List categories with amount per period
// @Description List categories with amount per period
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period path string true "period"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param filter query string false "Category filter"
// @Param order_by query string false "Sort field" example(name:asc,created_at:desc)
// @Success 200 {object} ListCategoryAmountPerPeriodResponse "List of categories with amount per period"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /categories/{period} [get]
func (api *API) ListCategoryAmountPerPeriod(ctx *gin.Context) {
	cmd := &ListAmountPerPeriodOptions{
		UseCases: api.categoriesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
