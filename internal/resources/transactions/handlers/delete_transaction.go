package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type DeleteOptions struct {
	Ctx      *gin.Context
	UseCases usecases.TransactionsUseCases

	UserID string
	ID     string
}

func (o *DeleteOptions) Complete(ctx *gin.Context) error {
	o.Ctx = ctx
	o.UserID = ctx.GetString("user_id")
	o.ID = ctx.Param("transaction_id")

	return nil
}

func (o *DeleteOptions) Validate() error {
	return nil
}

func (o *DeleteOptions) Run() error {
	err := o.UseCases.DeleteTransactionById(o.ID)
	if err != nil {
		return err
	}

	o.Ctx.Status(http.StatusNoContent)
	return nil
}

// @Summary Delete Transaction By ID
// @Description Delete a transaction and all entries related by the ID of the transaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param transaction_id path string true "transaction ID"
// @Success 204 "Transaction deleted"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 404 {object} utils.HTTPError "Not found"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/{transaction_id} [delete]
func (api *API) DeleteTransaction(ctx *gin.Context) {
	cmd := &DeleteOptions{
		UseCases: api.transactionsUseCases,
	}
	utils.RunCommand(ctx, cmd)
}
