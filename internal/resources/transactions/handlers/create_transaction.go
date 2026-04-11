package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CreateTransactionOptions struct {
	Ctx      *gin.Context
	UseCases usecases.TransactionsUseCases

	UserID string
	Body   transactions.CreateTransactionRequest
}

func (o *CreateTransactionOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")

	if err := ctx.ShouldBindJSON(&o.Body); err != nil {
		return utils.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *CreateTransactionOptions) Validate() error {
	return nil
}

func (o *CreateTransactionOptions) Run() error {
	var entriesDTO []usecases.CreateEntryDTO
	if o.Body.Entries != nil {
		entries := make([]usecases.CreateEntryDTO, len(o.Body.Entries))
		for i, e := range o.Body.Entries {
			t, _ := time.Parse("2006-01-02", e.ReferenceDate)
			entries[i] = usecases.CreateEntryDTO{
				Amount:        e.Amount,
				ReferenceDate: t,
			}
		}
		entriesDTO = entries
	}

	transaction, err := o.UseCases.CreateTransaction(usecases.CreateTransactionDTO{
		UserID:     o.UserID,
		Name:       o.Body.Name,
		CategoryID: utils.OptionalNullable[string]{Set: o.Body.CategoryID != nil, Value: o.Body.CategoryID},
		Note:       utils.OptionalNullable[string]{Set: o.Body.Note != nil, Value: o.Body.Note},
		Type:       o.Body.Type,
		Entries:    entriesDTO,
	})

	if err != nil {
		return err
	}

	o.Ctx.JSON(http.StatusCreated, transactions.CreateTransactionResponse{
		Data: transactions.CreateTransactionResponseData{
			Transaction: transaction,
		},
	})
	return nil
}

// @Summary Create a transaction
// @Description Create a transaction with all of it entries
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body transactions.CreateTransactionRequest true "Transaction payload"
// @Success 201 {object} transactions.CreateTransactionResponse "Installment updated"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /api/v1/transactions [post]
func (api *API) CreateTransaction(ctx *gin.Context) {
	cmd := &CreateTransactionOptions{
		UseCases: api.transactionsUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
