package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

// @gen_swagger_filter
var SummaryFilterConfig = querybuilder.ParseConfig{
	AllowedFields: map[string]querybuilder.FieldConfig{
		"period":        {AllowedOperators: []string{"eq", "in", "gte", "lte"}},
		"total_expense": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"total_income":  {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
		"total_balance": {AllowedOperators: []string{"eq", "gt", "gte", "lt", "lte"}},
	},
	AllowedSortFields: []string{"period", "total_expense", "total_income", "total_balance"},
}

type SummaryOptions struct {
	Ctx      *gin.Context
	UseCases usecases.TransactionsUseCases

	UserID  string
	Builder *querybuilder.Builder
}

func (o *SummaryOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")
	o.Builder = ctx.MustGet("query_builder").(*querybuilder.Builder).And("user_id", "eq", o.UserID)

	return nil
}

func (o *SummaryOptions) Validate() error {
	gte := o.Builder.HasAndFieldOperator("period", "gte")
	if len(gte) != 1 {
		return utils.NewHTTPError(http.StatusBadRequest, "exactly one 'period gte' filter is required")
	}
	gteVal, ok := gte[0].Value.(string)
	if !ok || len(gteVal) != 6 {
		return utils.NewHTTPError(http.StatusBadRequest, "period gte must be in YYYYMM format")
	}
	t1, err := time.Parse("200601", gteVal)
	if err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, "period gte is not a valid date (YYYYMM)")
	}

	lte := o.Builder.HasAndFieldOperator("period", "lte")
	if len(lte) != 1 {
		return utils.NewHTTPError(http.StatusBadRequest, "exactly one 'period lte' filter is required")
	}
	lteVal, ok := lte[0].Value.(string)
	if !ok || len(lteVal) != 6 {
		return utils.NewHTTPError(http.StatusBadRequest, "period lte must be in YYYYMM format")
	}
	t2, err := time.Parse("200601", lteVal)
	if err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, "period lte is not a valid date (YYYYMM)")
	}

	if t1.After(t2) {
		return utils.NewHTTPError(http.StatusBadRequest, "period gte must be lower than or equal to period lte")
	}

	months := (t2.Year()-t1.Year())*12 + int(t2.Month()) - int(t1.Month()) + 1
	if months > 12 {
		return utils.NewHTTPError(http.StatusBadRequest, "period range cannot exceed one year (12 months)")
	}

	return nil
}

func (o *SummaryOptions) Run() error {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(o.Ctx.Request.Context(), "TransactionsHandler.Summary")
	defer span.End()
	span.SetAttributes(attribute.String("user.id", o.UserID))

	reqCtx := querybuilder.WithBuilder(tCtx, o.Builder)
	summaryDTO, err := o.UseCases.Summary(reqCtx)

	if err != nil {
		span.RecordError(err)
		return err
	}

	summary := make([]MonthlySummaryResource, len(summaryDTO))
	for i, s := range summaryDTO {
		summary[i] = MonthlySummaryResource{
			Period:  s.Period,
			Income:  s.TotalIncome,
			Expense: s.TotalExpense,
			Balance: s.TotalBalance,
		}
	}

	o.Ctx.JSON(http.StatusOK, utils.ResponseData[SummaryResponseData]{
		Data: SummaryResponseData{
			Summary: summary,
		},
	})
	return nil
}

// @Summary Get financial summary grouped by month
// @Description Returns total income, expense and balance for each month in the specified period range.
// @Description Note: Only periods with existing transactions/entries will be returned.
// @ID v1GetSummary
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filter query string false "Filter expression. \n- Allowed fields & ops:\n  - period: eq, in, gte, lte\n  - total_balance: eq, gt, gte, lt, lte\n  - total_expense: eq, gt, gte, lt, lte\n  - total_income: eq, gt, gte, lt, lte\n"
// @Param order_by query string false "Sort field. \n- Allowed: period, total_expense, total_income, total_balance" example(period:asc)
// @Success 200 {object} utils.ResponseData[SummaryResponseData] "Summary data"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/transactions/summary [get]
func (api *API) Summary(ctx *gin.Context) {
	cmd := &SummaryOptions{
		UseCases: api.transactionsUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
