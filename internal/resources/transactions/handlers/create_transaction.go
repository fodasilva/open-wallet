package handlers

import (
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
)

type CreateTransactionOptions struct {
	W        http.ResponseWriter
	R        *http.Request
	UseCases usecases.TransactionsUseCases

	UserID string
	Body   CreateTransactionRequest
}

func (o *CreateTransactionOptions) Complete(w http.ResponseWriter, r *http.Request) error {
	o.W = w
	o.R = r
	o.UserID = util.GetString(r.Context(), util.ContextKeyUserID)

	if err := httputil.BindJSON(r, &o.Body); err != nil {
		return httputil.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (o *CreateTransactionOptions) Validate() error {
	if len(o.Body.Name) == 0 {
		return httputil.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	if string(o.Body.Type) == "" {
		return httputil.NewHTTPError(http.StatusBadRequest, "type is required")
	}
	if len(o.Body.Entries) == 0 {
		return httputil.NewHTTPError(http.StatusBadRequest, "entries are required")
	}
	for _, e := range o.Body.Entries {
		if e.ReferenceDate == "" {
			return httputil.NewHTTPError(http.StatusBadRequest, "reference_date is required for all entries")
		}
	}
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

	transaction, err := o.UseCases.CreateTransaction(o.R.Context(), usecases.CreateTransactionDTO{
		UserID:     o.UserID,
		Name:       o.Body.Name,
		CategoryID: util.OptionalNullable[string]{Set: o.Body.CategoryID != nil, Value: o.Body.CategoryID},
		Note:       util.OptionalNullable[string]{Set: o.Body.Note != nil, Value: o.Body.Note},
		Type:       o.Body.Type,
		Entries:    entriesDTO,
	})

	if err != nil {
		return err
	}

	httputil.JSON(o.W, http.StatusCreated, util.ResponseData[CreateTransactionResponseData]{
		Data: CreateTransactionResponseData{
			Transaction: MapTransactionResource(transaction),
		},
	})
	return nil
}

// @Summary Create a transaction
// @Description Create a transaction with all of it entries
// @ID v1CreateTransaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateTransactionRequest true "Transaction payload"
// @Success 201 {object} util.ResponseData[CreateTransactionResponseData] "Installment updated"
// @Failure 400 {object} httputil.HTTPError "Bad request"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Internal server error"
// @Failure 503 {string} string "Service Unavailable"
// @Router /api/v1/transactions [post]
func (api *API) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	cmd := &CreateTransactionOptions{
		UseCases: api.transactionsUseCases,
	}
	util.RunCommand(w, r, cmd)
}
