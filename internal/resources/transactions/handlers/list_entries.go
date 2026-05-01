package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

// @gen_swagger_filter
var EntriesFilterConfig = querybuilder.ParseConfig{
	AllowedFields: map[string]querybuilder.FieldConfig{
		"category_id":    {AllowedOperators: []string{"eq", "in"}},
		"type":           {AllowedOperators: []string{"eq", "in"}},
		"reference_date": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"amount":         {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"id":             {AllowedOperators: []string{"eq", "in"}},
		"user_id":        {AllowedOperators: []string{"eq", "in"}},
		"period":         {AllowedOperators: []string{"eq", "in", "gte", "lte"}},
		"created_at":     {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"transaction_id": {AllowedOperators: []string{"eq", "in"}},
	},
	AllowedSortFields: []string{"reference_date", "amount", "id", "created_at"},
}

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

	countCtx := querybuilder.WithBuilder(tCtx, querybuilder.ForCount(o.Builder))
	count, err := o.UseCases.CountEntries(countCtx)

	if err != nil {
		span.RecordError(err)
		return err
	}

	entriesResource := make([]EntryResource, len(entries))
	for i, e := range entries {
		entriesResource[i] = MapEntryResource(e)
	}

	o.Ctx.JSON(http.StatusOK, utils.PaginatedResponse[ListEntriesResponseData]{
		Data: ListEntriesResponseData{
			Entries: entriesResource,
		},
		Query: querybuilder.BuildMetadata(o.Page, o.PerPage, count),
	})
	return nil
}

// @Summary List entries
// @Description List a detailed view of entries joined with transactions for a given period
// @ID v1ListEntries
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param filter query string false "Filter expression. \n- Allowed fields & ops:\n  - amount: eq, gt, gte, lt, lte\n  - category_id: eq, in\n  - created_at: eq, gt, gte, lt, lte\n  - id: eq, in\n  - period: eq, in, gte, lte\n  - reference_date: eq, gt, gte, lt, lte\n  - transaction_id: eq, in\n  - type: eq, in\n  - user_id: eq, in\n"
// @Param order_by query string false "Sort field. \n- Allowed: reference_date, amount, id, created_at" example(reference_date:asc)
// @Success 200 {object} utils.PaginatedResponse[ListEntriesResponseData] "List of entries"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/transactions/entries [get]
func (api *API) ListEntries(ctx *gin.Context) {
	cmd := &ListEntriesOptions{
		UseCases: api.transactionsUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
