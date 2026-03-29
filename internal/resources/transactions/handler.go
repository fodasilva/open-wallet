package transactions

import (
	"database/sql"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/gin-gonic/gin"
)

type API struct {
	transactionsUseCase TransactionsUseCase
}

func NewHandler(db *sql.DB) *API {
	return &API{
		transactionsUseCase: NewTransactionsUseCase(NewTransactionsRepo(db),
			categories.NewCategoriesUseCase(repository.NewCategoriesRepo(), db),
			db),
	}
}

// @Summary List entries
// @Description List a detailed view of entries joined with transactions for a given period
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param order_by query string false "Sort field" example(name:asc,created_at:desc)
// @Success 200 {object} ListEntriesResponse "List of entries"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/entries [get]
func (api *API) ListEntries(ctx *gin.Context) {
	tracer := otel.Tracer("handler")
	tCtx, span := tracer.Start(ctx.Request.Context(), "TransactionsHandler.ListEntries")
	defer span.End()

	userID := ctx.GetString("user_id")
	span.SetAttributes(attribute.String("user.id", userID))
	page := ctx.GetInt("page")
	perPage := ctx.GetInt("per_page")
	queryOpts := ctx.MustGet("query_opts").(*utils.QueryOptsBuilder).And("user_id", "eq", userID)

	entries, err := api.transactionsUseCase.ListViewEntries(tCtx, queryOpts)

	if err != nil {
		span.RecordError(err)
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	count, err := api.transactionsUseCase.CountViewEntries(tCtx, utils.QueryOpts().
		And("user_id", "eq", userID))

	if err != nil {
		span.RecordError(err)
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	nextPage := len(entries) > perPage

	if nextPage {
		entries = entries[:len(entries)-1]
	}

	totalPages := (count + perPage - 1) / perPage

	ctx.JSON(http.StatusOK, ListEntriesResponse{
		Data: ListEntriesResponseData{
			Entries: entries,
		},
		Query: utils.QueryMeta{
			Page:       page,
			PerPage:    perPage,
			NextPage:   nextPage,
			TotalPages: totalPages,
			TotalItems: count,
		},
	})
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
	id := ctx.Param("transaction_id")

	err := api.transactionsUseCase.DeleteTransactionById(id)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Create a transaction
// @Description Create a transaction with all of it entries
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateTransactionRequest true "Transaction payload"
// @Success 200 {object} CreateTransactionResponse "Installment updated"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions [post]
func (api *API) CreateTransaction(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	var body CreateTransactionRequest

	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	var entriesDTO []CreateEntryDTO
	if body.Entries != nil {
		entries := make([]CreateEntryDTO, len(body.Entries))
		for i, entry := range body.Entries {
			entries[i] = CreateEntryDTO{
				Amount:        entry.Amount,
				ReferenceDate: entry.ReferenceDate,
			}
		}
		entriesDTO = entries
	}

	transaction, err := api.transactionsUseCase.CreateTransaction(CreateTransactionDTO{
		UserID:     userID,
		Name:       body.Name,
		CategoryID: body.CategoryID,
		Note:       body.Note,
		Type:       body.Type,
		Entries:    entriesDTO,
	})

	if err != nil {
		apiErr := utils.GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, CreateTransactionResponse{
		Data: CreateTransactionResponseData{
			Transaction: transaction,
		},
	})
}

// @Summary Update a transaction
// @Description Update a transaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param transaction_id path string true "transaction ID"
// @Param body body UpdateTransactionRequest true "Installment payload"
// @Success 200 {object} UpdateTransactionResponse "Installment updated"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/{transaction_id} [patch]
func (api *API) UpdateTransaction(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	transactionID := ctx.Param("transaction_id")
	var body UpdateTransactionRequest

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		ctx.Abort()
		return
	}

	var entriesDTO *[]UpdateEntryDTO
	if body.Entries != nil {
		entries := make([]UpdateEntryDTO, len(*body.Entries))
		for i, entry := range *body.Entries {
			entries[i] = UpdateEntryDTO{
				Amount:        entry.Amount,
				ReferenceDate: entry.ReferenceDate,
			}
		}
		entriesDTO = &entries
	}

	transaction, err := api.transactionsUseCase.UpdateTransaction(transactionID, userID, UpdateTransactionDTO{
		Update:     body.Update,
		Name:       body.Name,
		Note:       body.Note,
		CategoryID: body.CategoryID,
		Entries:    entriesDTO,
	})

	if err != nil {
		apiErr := utils.GetApiErr(err)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, UpdateTransactionResponse{
		Data: UpdateTransactionResponseData{
			Transaction: transaction,
		},
	})
}
