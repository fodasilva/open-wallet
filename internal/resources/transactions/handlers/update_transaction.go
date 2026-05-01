package handlers

import (
	"net/http"
	"slices"
	"time"

	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type UpdateTransactionOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.TransactionsUseCases

	ID         string
	UserID     string
	PassedKeys []string
	Body       UpdateTransactionRequest
}

func (o *UpdateTransactionOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.ID = r.PathValue("transaction_id")
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)

	passedKeys, err := util.GetJSONKeys(r)
	if err != nil {
		return util.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(passedKeys) == 0 {
		return util.NewHTTPError(http.StatusBadRequest, "At least one field must be provided for update")
	}
	o.PassedKeys = passedKeys

	if err := httputil.BindJSON(r, &o.Body); err != nil {
		return util.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *UpdateTransactionOptions) Validate() error {
	if slices.Contains(o.PassedKeys, "name") {
		if o.Body.Name == nil {
			return util.NewHTTPError(http.StatusBadRequest, "name cannot be null")
		}
		if len(*o.Body.Name) == 0 {
			return util.NewHTTPError(http.StatusBadRequest, "name cannot be empty")
		}
	}
	if slices.Contains(o.PassedKeys, "entries") {
		if o.Body.Entries == nil {
			return util.NewHTTPError(http.StatusBadRequest, "entries cannot be null")
		}
		if len(*o.Body.Entries) == 0 {
			return util.NewHTTPError(http.StatusBadRequest, "entries cannot be empty")
		}
		for _, e := range *o.Body.Entries {
			if e.ReferenceDate == "" {
				return util.NewHTTPError(http.StatusBadRequest, "reference_date is required for all entries")
			}
		}
	}
	return nil
}

func (o *UpdateTransactionOptions) Run() error {
	var payload usecases.UpdateTransactionDTO

	if slices.Contains(o.PassedKeys, "name") {
		payload.Name = util.OptionalNullable[string]{Set: true, Value: o.Body.Name}
	}

	if slices.Contains(o.PassedKeys, "note") {
		payload.Note = util.OptionalNullable[string]{Set: true, Value: o.Body.Note}
	}

	if slices.Contains(o.PassedKeys, "category_id") {
		payload.CategoryID = util.OptionalNullable[string]{Set: true, Value: o.Body.CategoryID}
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
		payload.Entries = util.OptionalNullable[[]usecases.UpdateEntryDTO]{Set: true, Value: &updatedEntries}
	}

	transaction, err := o.UseCases.UpdateTransaction(o.R.Context(), o.ID, o.UserID, payload)

	if err != nil {
		return err
	}

	httputil.JSON(o.W, http.StatusOK, util.ResponseData[UpdateTransactionResponseData]{
		Data: UpdateTransactionResponseData{
			Transaction: MapTransactionResource(transaction),
		},
	})
	return nil
}

// @Summary Update a transaction
// @Description Update a transaction
// @ID v1UpdateTransaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param transaction_id path string true "transaction ID"
// @Param body body UpdateTransactionRequest true "Installment payload"
// @Success 200 {object} util.ResponseData[UpdateTransactionResponseData] "Installment updated"
// @Failure 400 {object} util.HTTPError "Bad request"
// @Failure 401 {object} util.HTTPError "Unauthorized"
// @Failure 500 {object} util.HTTPError "Internal server error"
// @Router /api/v1/transactions/{transaction_id} [patch]
func (api *API) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	cmd := &UpdateTransactionOptions{
		UseCases: api.transactionsUseCases,
	}
	util.RunCommand(w, r, cmd)
}
