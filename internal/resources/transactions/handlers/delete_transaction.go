package handlers

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
)

type DeleteOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.TransactionsUseCases

	UserID string
	ID     string
}

func (o *DeleteOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)
	o.ID = r.PathValue("transaction_id")

	return nil
}

func (o *DeleteOptions) Validate() error {
	return nil
}

func (o *DeleteOptions) Run() error {
	err := o.UseCases.DeleteTransactionById(o.R.Context(), o.ID, o.UserID)
	if err != nil {
		return err
	}

	o.W.WriteHeader(http.StatusNoContent)
	return nil
}

// @Summary Delete Transaction By ID
// @Description Delete a transaction and all entries related by the ID of the transaction
// @ID v1DeleteTransaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param transaction_id path string true "transaction ID"
// @Success 204 "Transaction deleted"
// @Failure 401 {object} util.HTTPError "Unauthorized"
// @Failure 404 {object} util.HTTPError "Not found"
// @Failure 500 {object} util.HTTPError "Internal server error"
// @Router /api/v1/transactions/{transaction_id} [delete]
func (api *API) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	cmd := &DeleteOptions{
		UseCases: api.transactionsUseCases,
	}
	util.RunCommand(w, r, cmd)
}
