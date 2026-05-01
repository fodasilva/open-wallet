package handlers

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

// @gen_swagger_filter
var TransactionsFilterConfig = querybuilder.ParseConfig{
	AllowedFields: map[string]querybuilder.FieldConfig{
		"category_id":    {AllowedOperators: []string{"eq", "in"}},
		"type":           {AllowedOperators: []string{"eq", "in"}},
		"reference_date": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"amount":         {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"id":             {AllowedOperators: []string{"eq", "in"}},
		"user_id":        {AllowedOperators: []string{"eq", "in"}},
		"period":         {AllowedOperators: []string{"eq", "in", "gte", "lte"}},
		"created_at":     {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
	},
	AllowedSortFields: []string{"reference_date", "amount", "id", "created_at"},
}

type ListEntriesOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.TransactionsUseCases

	UserID  string
	Page    int
	PerPage int
	Builder *querybuilder.Builder
}

func (o *ListEntriesOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)
	o.Page = util.GetInt(r.Context(), util.ContextKeyPage)
	o.PerPage = util.GetInt(r.Context(), util.ContextKeyPerPage)
	o.Builder = querybuilder.Get(r.Context()).And("user_id", "eq", o.UserID)

	return nil
}

func (o *ListEntriesOptions) Validate() error {
	return nil
}

func (o *ListEntriesOptions) Run() error {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(o.R.Context(), "TransactionsHandler.ListEntries")
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

	httputil.JSON(o.W, http.StatusOK, util.PaginatedResponse[ListEntriesResponseData]{
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
// @Param filter query string false "Filter expression. \n- Allowed fields & ops:\n  - amount: eq, gt, gte, lt, lte\n  - category_id: eq, in\n  - created_at: eq, gt, gte, lt, lte\n  - id: eq, in\n  - period: eq, in, gte, lte\n  - reference_date: eq, gt, gte, lt, lte\n  - type: eq, in\n  - user_id: eq, in\n"
// @Param order_by query string false "Sort field. \n- Allowed: reference_date, amount, id, created_at" example(reference_date:asc)
// @Success 200 {object} util.PaginatedResponse[ListEntriesResponseData] "List of entries"
// @Failure 401 {object} util.HTTPError "Unauthorized"
// @Failure 500 {object} util.HTTPError "Internal server error"
// @Router /api/v1/transactions/entries [get]
func (api *API) ListEntries(w http.ResponseWriter, r *http.Request) {
	cmd := &ListEntriesOptions{
		UseCases: api.transactionsUseCases,
	}
	util.RunCommand(w, r, cmd)
}
