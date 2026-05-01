package handlers

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

// @gen_swagger_filter
var RecurrencesFilterConfig = querybuilder.ParseConfig{
	AllowedFields: map[string]querybuilder.FieldConfig{
		"id":          {AllowedOperators: []string{"eq", "in"}},
		"category_id": {AllowedOperators: []string{"eq", "in"}},
		"name":        {AllowedOperators: []string{"eq", "like", "in"}},
		"user_id":     {AllowedOperators: []string{"eq", "in"}},
		"created_at":  {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"amount":      {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
	},
	AllowedSortFields: []string{"name", "created_at", "id"},
}

type ListOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.RecurrencesUseCases

	UserID  string
	Page    int
	PerPage int
	Builder *querybuilder.Builder
}

func (o *ListOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)
	o.Page = util.GetInt(r.Context(), util.ContextKeyPage)
	o.PerPage = util.GetInt(r.Context(), util.ContextKeyPerPage)
	o.Builder = querybuilder.Get(r.Context()).And("user_id", "eq", o.UserID)

	return nil
}

func (o *ListOptions) Validate() error {
	return nil
}

func (o *ListOptions) Run() error {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(o.R.Context(), "RecurrencesHandler.List")
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

	recurrencesResource := make([]RecurrenceResource, len(items))
	for i, r := range items {
		recurrencesResource[i] = MapRecurrenceResource(r)
	}

	httputil.JSON(o.W, http.StatusOK, util.PaginatedResponse[ListRecurrencesResponseData]{
		Data: ListRecurrencesResponseData{
			Recurrences: recurrencesResource,
		},
		Query: querybuilder.BuildMetadata(o.Page, o.PerPage, count),
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
// @Param filter query string false "Filter expression. \n- Allowed fields & ops:\n  - amount: eq, gt, gte, lt, lte\n  - category_id: eq, in\n  - created_at: eq, gt, gte, lt, lte\n  - id: eq, in\n  - name: eq, like, in\n  - user_id: eq, in\n"
// @Param order_by query string false "Sort field. \n- Allowed: name, created_at, id" example(name:asc)
// @Success 200 {object} util.PaginatedResponse[ListRecurrencesResponseData] "List of recurrences"
// @Failure 401 {object} util.HTTPError "Unauthorized"
// @Failure 500 {object} util.HTTPError "Internal server error"
// @Router /api/v1/recurrences [get]
func (api *API) List(w http.ResponseWriter, r *http.Request) {
	cmd := &ListOptions{
		UseCases: api.recurrencesUseCases,
	}
	util.RunCommand(w, r, cmd)
}
