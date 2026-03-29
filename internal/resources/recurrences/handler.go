package recurrences

import (
	"database/sql"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/gin-gonic/gin"
)

type API struct {
	recurrencesUseCase RecurrencesUseCase
}

func NewHandler(db *sql.DB) *API {
	catsRepo := repository.NewCategoriesRepo()
	catsUseCase := categories.NewCategoriesUseCase(catsRepo, db)
	txsRepo := transactions.NewTransactionsRepo(db)
	txsUseCase := transactions.NewTransactionsUseCase(txsRepo, catsUseCase, db)

	return &API{
		recurrencesUseCase: NewRecurrencesUseCase(NewRecurrencesRepo(),
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
	var body CreateRecurrenceRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	rec, err := api.recurrencesUseCase.Create(CreateRecurrenceDTO{
		UserID:      userID,
		Name:        body.Name,
		CategoryID:  body.CategoryID,
		Note:        body.Note,
		Amount:      body.Amount,
		DayOfMonth:  body.DayOfMonth,
		StartPeriod: body.StartPeriod,
		EndPeriod:   body.EndPeriod,
	})

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
	var body UpdateRecurrenceRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		ctx.Abort()
		return
	}

	rec, err := api.recurrencesUseCase.Update(id, userID, UpdateRecurrenceDTO{
		Update:      body.Update,
		Name:        body.Name,
		CategoryID:  body.CategoryID,
		Note:        body.Note,
		Amount:      body.Amount,
		DayOfMonth:  body.DayOfMonth,
		StartPeriod: body.StartPeriod,
		EndPeriod:   body.EndPeriod,
	})

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
