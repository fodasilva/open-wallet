package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

type ListEntriesOptions struct {
	Ctx      *gin.Context
	UseCases usecases.TransactionsUseCases

	UserID  string
	Page    int
	PerPage int
	Builder *querybuilder.Builder
}

func (o *ListEntriesOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")
	o.Page = ctx.GetInt("page")
	o.PerPage = ctx.GetInt("per_page")
	o.Builder = ctx.MustGet("query_builder").(*querybuilder.Builder).And("user_id", "eq", o.UserID)

	return nil
}

func (o *ListEntriesOptions) Validate() error {
	return nil
}

func (o *ListEntriesOptions) Run() error {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(o.Ctx.Request.Context(), "TransactionsHandler.ListEntries")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", o.UserID))

	reqCtx := querybuilder.WithBuilder(tCtx, o.Builder)
	entries, err := o.UseCases.ListEntries(reqCtx)

	if err != nil {
		span.RecordError(err)
		return err
	}

	countCtx := querybuilder.WithBuilder(tCtx, querybuilder.New().And("user_id", "eq", o.UserID))
	count, err := o.UseCases.CountEntries(countCtx)

	if err != nil {
		span.RecordError(err)
		return err
	}

	nextPage := len(entries) > o.PerPage

	if nextPage {
		entries = entries[:len(entries)-1]
	}

	totalPages := (count + o.PerPage - 1) / o.PerPage

	o.Ctx.JSON(http.StatusOK, transactions.ListEntriesResponse{
		Data: transactions.ListEntriesResponseData{
			Entries: entries,
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

// @Summary List entries
// @Description List a detailed view of entries joined with transactions for a given period
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param order_by query string false "Sort field" example(name:asc,created_at:desc)
// @Success 200 {object} transactions.ListEntriesResponse "List of entries"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/transactions/entries [get]
func (api *API) ListEntries(ctx *gin.Context) {
	cmd := &ListEntriesOptions{
		UseCases: api.transactionsUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
