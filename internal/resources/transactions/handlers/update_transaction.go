package handlers

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type UpdateTransactionOptions struct {
	Ctx      *gin.Context
	UseCases usecases.TransactionsUseCases

	ID         string
	UserID     string
	PassedKeys []string
	Body       transactions.UpdateTransactionRequest
}

func (o *UpdateTransactionOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.ID = ctx.Param("transaction_id")
	o.UserID = ctx.GetString("user_id")

	passedKeys, err := utils.GetJSONKeys(ctx)
	if err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(passedKeys) == 0 {
		return utils.NewHTTPError(http.StatusBadRequest, "At least one field must be provided for update")
	}
	o.PassedKeys = passedKeys

	if err := ctx.ShouldBindJSON(&o.Body); err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *UpdateTransactionOptions) Validate() error {
	if slices.Contains(o.PassedKeys, "name") && o.Body.Name == nil {
		return utils.NewHTTPError(http.StatusBadRequest, "name cannot be null")
	}
	if slices.Contains(o.PassedKeys, "entries") && o.Body.Entries == nil {
		return utils.NewHTTPError(http.StatusBadRequest, "entries cannot be null")
	}
	return nil
}

func (o *UpdateTransactionOptions) Run() error {
	var payload usecases.UpdateTransactionDTO

	if slices.Contains(o.PassedKeys, "name") {
		payload.Name = utils.OptionalNullable[string]{Set: true, Value: o.Body.Name}
	}

	if slices.Contains(o.PassedKeys, "note") {
		payload.Note = utils.OptionalNullable[string]{Set: true, Value: o.Body.Note}
	}

	if slices.Contains(o.PassedKeys, "category_id") {
		payload.CategoryID = utils.OptionalNullable[string]{Set: true, Value: o.Body.CategoryID}
	}

	if slices.Contains(o.PassedKeys, "entries") {
		updatedEntries := make([]usecases.UpdateEntryDTO, len(*o.Body.Entries))
		for i, e := range *o.Body.Entries {
			t, _ := time.Parse("2006-01-02", e.ReferenceDate)
			updatedEntries[i] = usecases.UpdateEntryDTO{
				Amount:        e.Amount,
				ReferenceDate: t,
			}
		}
		payload.Entries = utils.OptionalNullable[[]usecases.UpdateEntryDTO]{Set: true, Value: &updatedEntries}
	}

	transaction, err := o.UseCases.UpdateTransaction(o.Ctx.Request.Context(), o.ID, o.UserID, payload)

	if err != nil {
		return err
	}

	o.Ctx.JSON(http.StatusOK, transactions.UpdateTransactionResponse{
		Data: transactions.UpdateTransactionResponseData{
			Transaction: transaction,
		},
	})
	return nil
}

// @Summary Update a transaction
// @Description Update a transaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param transaction_id path string true "transaction ID"
// @Param body body transactions.UpdateTransactionRequest true "Installment payload"
// @Success 200 {object} transactions.UpdateTransactionResponse "Installment updated"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/transactions/{transaction_id} [patch]
func (api *API) UpdateTransaction(ctx *gin.Context) {
	cmd := &UpdateTransactionOptions{
		UseCases: api.transactionsUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
