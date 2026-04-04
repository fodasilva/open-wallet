package recurrences

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"slices"

	"github.com/felipe1496/open-wallet/internal/resources/categories"
	categoriesRepository "github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	transactionsRepository "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type API struct {
	recurrencesUseCase RecurrencesUseCase
}

func NewHandler(db *sql.DB) *API {
	catsRepo := categoriesRepository.NewCategoriesRepo()
	catsUseCase := categories.NewCategoriesUseCase(catsRepo, db)
	txsRepo := transactionsRepository.NewTransactionsRepo()
	entriesRepo := transactionsRepository.NewEntriesRepo()
	txsUseCase := transactions.NewTransactionsUseCase(txsRepo, entriesRepo, catsUseCase, db)

	return &API{
		recurrencesUseCase: NewRecurrencesUseCase(repository.NewRecurrencesRepo(),
			catsUseCase,
			txsUseCase,
			db),
	}
}

// @Summary Create a recurrence
// @Description Create a new recurrence template
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateRecurrenceRequest true "Recurrence payload"
// @Success 201 {object} CreateRecurrenceResponse "Recurrence created"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences [post]
func (api *API) Create(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	passedKeys, err := utils.GetJSONKeys(ctx)
	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	var body CreateRecurrenceRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	payload := repository.CreateRecurrenceDTO{
		UserID:      userID,
		Name:        body.Name,
		CategoryID:  utils.OptionalNullable[string]{Set: slices.Contains(passedKeys, "category_id"), Value: body.CategoryID},
		Note:        utils.OptionalNullable[string]{Set: slices.Contains(passedKeys, "note"), Value: body.Note},
		Amount:      body.Amount,
		DayOfMonth:  body.DayOfMonth,
		StartPeriod: body.StartPeriod,
		EndPeriod:   utils.OptionalNullable[string]{Set: slices.Contains(passedKeys, "end_period"), Value: body.EndPeriod},
	}

	rec, err := api.recurrencesUseCase.Create(payload)

	if err != nil {
		apiErr := utils.GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, CreateRecurrenceResponse{
		Data: CreateRecurrenceResponseData{
			Recurrence: rec,
		},
	})
}

// @Summary List recurrences
// @Description List user recurrences
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} ListRecurrencesResponse "List of recurrences"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences [get]
func (api *API) List(ctx *gin.Context) {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(ctx.Request.Context(), "RecurrencesHandler.List")
	defer span.End()

	userID := ctx.GetString("user_id")
	span.SetAttributes(attribute.String("user.id", userID))
	page := ctx.GetInt("page")
	perPage := ctx.GetInt("per_page")
	queryOpts := ctx.MustGet("query_opts").(*utils.QueryOptsBuilder).And("user_id", "eq", userID)

	items, err := api.recurrencesUseCase.List(tCtx, queryOpts)
	if err != nil {
		span.RecordError(err)
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	count, err := api.recurrencesUseCase.Count(tCtx, utils.QueryOpts().And("user_id", "eq", userID))
	if err != nil {
		span.RecordError(err)
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	nextPage := len(items) > perPage
	if nextPage {
		items = items[:len(items)-1]
	}
	totalPages := (count + perPage - 1) / perPage

	ctx.JSON(http.StatusOK, ListRecurrencesResponse{
		Data: ListRecurrencesResponseData{
			Recurrences: items,
		},
		Query: utils.QueryMeta{
			Page:       page,
			PerPage:    perPage,
			NextPage:   nextPage,
			TotalPages: totalPages,
			TotalItems: count,
		},
	})
}

// @Summary Delete Recurrence By ID
// @Description Delete a recurrence template
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "recurrence ID"
// @Success 204 "Recurrence deleted"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 404 {object} utils.HTTPError "Not found"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences/{id} [delete]
func (api *API) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	scope := ctx.DefaultQuery("scope", "all")

	err := api.recurrencesUseCase.DeleteByID(id, scope)
	if err != nil {
		apiErr := utils.GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Update a recurrence
// @Description Update a recurrence template
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "recurrence ID"
// @Param body body UpdateRecurrenceRequest true "Recurrence payload"
// @Success 200 {object} UpdateRecurrenceResponse "Recurrence updated"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences/{id} [patch]
func (api *API) Update(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	id := ctx.Param("id")
	passedKeys, err := utils.GetJSONKeys(ctx)
	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	if len(passedKeys) == 0 {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, "At least one field must be provided for update")
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	var body UpdateRecurrenceRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		ctx.Abort()
		return
	}

	// Validate non-nullable fields
	nonNullable := map[string]interface{}{
		"name":         body.Name,
		"amount":       body.Amount,
		"day_of_month": body.DayOfMonth,
		"start_period": body.StartPeriod,
	}

	for _, key := range passedKeys {
		if val, ok := nonNullable[key]; ok && (val == nil || reflect.ValueOf(val).IsNil()) {
			apiErr := utils.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s cannot be null", key))
			ctx.JSON(apiErr.StatusCode, apiErr)
			return
		}
	}

	payload := repository.UpdateRecurrenceDTO{
		Name:        utils.OptionalNullable[string]{Set: slices.Contains(passedKeys, "name"), Value: body.Name},
		CategoryID:  utils.OptionalNullable[string]{Set: slices.Contains(passedKeys, "category_id"), Value: body.CategoryID},
		Note:        utils.OptionalNullable[string]{Set: slices.Contains(passedKeys, "note"), Value: body.Note},
		Amount:      utils.OptionalNullable[float64]{Set: slices.Contains(passedKeys, "amount"), Value: body.Amount},
		DayOfMonth:  utils.OptionalNullable[int]{Set: slices.Contains(passedKeys, "day_of_month"), Value: body.DayOfMonth},
		StartPeriod: utils.OptionalNullable[string]{Set: slices.Contains(passedKeys, "start_period"), Value: body.StartPeriod},
		EndPeriod:   utils.OptionalNullable[string]{Set: slices.Contains(passedKeys, "end_period"), Value: body.EndPeriod},
	}

	rec, err := api.recurrencesUseCase.Update(id, userID, payload)

	if err != nil {
		apiErr := utils.GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, UpdateRecurrenceResponse{
		Data: UpdateRecurrenceResponseData{
			Recurrence: rec,
		},
	})
}

// @Summary Prepare recurrences for a period
// @Description Generates entry records for all recurrence templates that don't already have one in the given period.
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period path string true "Period in YYYYMM format (e.g. 202603)"
// @Success 200 "Recurrences prepared"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /recurrences/{period} [post]
func (api *API) Prepare(ctx *gin.Context) {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(ctx.Request.Context(), "RecurrencesHandler.Prepare")
	defer span.End()

	userID := ctx.GetString("user_id")
	span.SetAttributes(attribute.String("user.id", userID))

	period := ctx.Param("period")
	if len(period) != 6 {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, "invalid period format. Expected YYYYMM.")
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	err := api.recurrencesUseCase.PrepareRecurrences(tCtx, userID, period)
	if err != nil {
		span.RecordError(err)
		apiErr := utils.GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
