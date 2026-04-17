package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

type ListOptions struct {
	Ctx      *gin.Context
	UseCases usecases.RecurrencesUseCases

	UserID  string
	Page    int
	PerPage int
	Builder *querybuilder.Builder
}

func (o *ListOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")
	o.Page = ctx.GetInt("page")
	o.PerPage = ctx.GetInt("per_page")
	o.Builder = ctx.MustGet("query_builder").(*querybuilder.Builder).And("user_id", "eq", o.UserID)

	return nil
}

func (o *ListOptions) Validate() error {
	return nil
}

func (o *ListOptions) Run() error {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(o.Ctx.Request.Context(), "RecurrencesHandler.List")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", o.UserID))

	reqCtx := querybuilder.WithBuilder(tCtx, o.Builder)
	items, err := o.UseCases.List(reqCtx)
	if err != nil {
		span.RecordError(err)
		return err
	}

	countCtx := querybuilder.WithBuilder(tCtx, querybuilder.New().And("user_id", "eq", o.UserID))
	count, err := o.UseCases.Count(countCtx)
	if err != nil {
		span.RecordError(err)
		return err
	}

	nextPage := len(items) > o.PerPage
	if nextPage {
		items = items[:len(items)-1]
	}
	totalPages := (count + o.PerPage - 1) / o.PerPage

	o.Ctx.JSON(http.StatusOK, ListRecurrencesResponse{
		Data: ListRecurrencesResponseData{
			Recurrences: items,
		},
		Query: utils.QueryMeta{
			Page:       o.Page,
			PerPage:    o.PerPage,
			NextPage:   nextPage,
			TotalPages: totalPages,
			TotalItems: count,
		},
	})

	return nil
}

// @Summary List recurrences
// @Description List user recurrences
// @ID v1ListRecurrences
// @Tags recurrences
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} ListRecurrencesResponse "List of recurrences"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/recurrences [get]
func (api *API) List(ctx *gin.Context) {
	cmd := &ListOptions{
		UseCases: api.recurrencesUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
