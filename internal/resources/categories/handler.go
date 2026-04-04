package categories

import (

	"net/http"
	"slices"
	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/gin-gonic/gin"
)

type API struct {
	categoriesUseCase CategoriesUseCase
}

func NewHandler(categoriesUseCase CategoriesUseCase) *API {
	return &API{
		categoriesUseCase: categoriesUseCase,
	}
}

// @Summary Create a category
// @Description Create a category
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateCategoryRequest true "Category payload"
// @Success 201 {object} CreateCategoryResponse "Category created"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /categories [post]
func (api *API) Create(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	var body CreateCategoryRequest

	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	category, err := api.categoriesUseCase.Create(repository.CreateCategoryDTO{
		UserID: userID,
		Name:   body.Name,
		Color:  body.Color,
	})

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusInternalServerError, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, CreateCategoryResponse{
		Data: CreateCategoryResponseData{
			Category: category,
		},
	})
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
// @Router /categories [get]
func (api *API) List(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	queryOpts := ctx.MustGet("query_opts").(*utils.QueryOptsBuilder).And("user_id", "eq", userID)
	nameFilter := ctx.Query("name")

	if nameFilter != "" {
		queryOpts.And("name", "like", nameFilter)
	}

	categories, err := api.categoriesUseCase.List(queryOpts)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusInternalServerError, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	count, err := api.categoriesUseCase.Count(utils.ForCount(queryOpts))
	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusInternalServerError, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	page := ctx.GetInt("page")
	perPage := ctx.GetInt("per_page")
	nextPage := len(categories) > perPage
	totalPages := (count + perPage - 1) / perPage

	if nextPage {
		categories = categories[:len(categories)-1]
	}

	ctx.JSON(http.StatusOK, ListCategoriesResponse{
		Data: ListCategoriesResponseData{
			Categories: categories,
		},
		Query: utils.QueryMeta{
			NextPage:   nextPage,
			Page:       page,
			PerPage:    perPage,
			TotalItems: count,
			TotalPages: totalPages,
		},
	})
}

// @Summary Delete Category By ID
// @Description Delete a category
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category_id path string true "category ID"
// @Success 204 "Category deleted"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 404 {object} utils.HTTPError "Not found"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /categories/{category_id} [delete]
func (api *API) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("category_id")

	err := api.categoriesUseCase.DeleteByID(id)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.Status(http.StatusNoContent)
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
	userID := ctx.GetString("user_id")
	period := ctx.Param("period")
	queryOpts := ctx.MustGet("query_opts").(*utils.QueryOptsBuilder).
		And("user_id", "eq", userID)

	categories, err := api.categoriesUseCase.ListCategoryAmountPerPeriod(period, queryOpts)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusInternalServerError, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	count, err := api.categoriesUseCase.CountCategoryAmountPerPeriod(period, utils.ForCount(queryOpts))
	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	page := ctx.GetInt("page")
	perPage := ctx.GetInt("per_page")
	nextPage := len(categories) > perPage
	totalPages := (count + perPage - 1) / perPage

	if nextPage {
		categories = categories[:len(categories)-1]
	}

	ctx.JSON(http.StatusOK, ListCategoryAmountPerPeriodResponse{
		Data: ListCategoryAmountPerPeriodResponseData{
			Categories: categories,
		},
		Query: utils.QueryMeta{
			NextPage:   nextPage,
			Page:       page,
			PerPage:    perPage,
			TotalItems: count,
			TotalPages: totalPages,
		},
	})
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
// @Router /categories/{category_id} [patch]
func (api *API) Update(ctx *gin.Context) {
	id := ctx.Param("category_id")
	passedKeys, err := utils.GetJSONKeys(ctx)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	if len(passedKeys) == 0 {
		apiErr := utils.NewHTTPError(
			http.StatusBadRequest,
			"At least one field must be provided for update",
		)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	var body UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	var payload repository.UpdateCategoryDTO

	if slices.Contains(passedKeys, "name") {
		if body.Name == nil {
			apiErr := utils.NewHTTPError(http.StatusBadRequest, "name cannot be null")
			ctx.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
			return
		}
		payload.Name = utils.NewValue(*body.Name)
	}

	if slices.Contains(passedKeys, "color") {
		if body.Color == nil {
			apiErr := utils.NewHTTPError(http.StatusBadRequest, "color cannot be null")
			ctx.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
			return
		}
		payload.Color = utils.NewValue(*body.Color)
	}

	category, err := api.categoriesUseCase.Update(id, payload)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, UpdateCategoryResponse{
		Data: UpdateCategoryResponseData{
			Category: category,
		},
	})
}
