package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *RecurrencesUseCasesImpl) Update(ctx context.Context, id string, userID string, payload repository.UpdateRecurrenceDTO) (repository.Recurrence, error) {
	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().
		And("id", "eq", id).
		And("user_id", "eq", userID))
	exists, err := uc.repo.Select(filterCtx, uc.db)

	if err != nil {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if recurrence exists")
	}

	if len(exists) == 0 {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusNotFound, "recurrence not found")
	}

	if err := uc.validateCategory(ctx, userID, payload.CategoryID); err != nil {
		return repository.Recurrence{}, err
	}

	updateFilterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", id))
	err = uc.repo.Update(updateFilterCtx, uc.db, payload)
	if err != nil {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update recurrence")
	}

	rec, err := uc.fetchRecurrence(ctx, id)
	if err != nil {
		return repository.Recurrence{}, err
	}

	if payload.Amount.Set && payload.Amount.Value != nil {
		if err := uc.syncLinkedTransactions(ctx, id, userID, *payload.Amount.Value); err != nil {
			return rec, err
		}
	}

	return rec, nil
}

func (uc *RecurrencesUseCasesImpl) validateCategory(ctx context.Context, userID string, categoryID utils.OptionalNullable[string]) error {
	if !categoryID.Set || categoryID.Value == nil {
		return nil
	}

	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().
		And("id", "eq", *categoryID.Value).
		And("user_id", "eq", userID))
	exists, err := uc.categoriesUseCase.List(filterCtx)
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
	}

	if len(exists) == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	return nil
}

func (uc *RecurrencesUseCasesImpl) fetchRecurrence(ctx context.Context, id string) (repository.Recurrence, error) {
	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", id))
	recs, err := uc.repo.Select(filterCtx, uc.db)
	if err != nil || len(recs) == 0 {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch recurrence")
	}
	return recs[0], nil
}

func (uc *RecurrencesUseCasesImpl) syncLinkedTransactions(ctx context.Context, id string, userID string, amount float64) error {
	txFilterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().
		And("user_id", "eq", userID).
		And("recurrence_id", "eq", id))
	txs, err := uc.transactionsUseCase.ListEntries(txFilterCtx)
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch linked transactions for sync")
	}

	if len(txs) == 0 {
		return nil
	}

	transactionID := txs[0].TransactionID
	updatedEntries := make([]usecases.UpdateEntryDTO, len(txs))
	for i, entry := range txs {
		updatedEntries[i] = usecases.UpdateEntryDTO{
			Amount:        amount,
			ReferenceDate: entry.ReferenceDate,
		}
	}

	_, err = uc.transactionsUseCase.UpdateTransaction(ctx, transactionID, userID, usecases.UpdateTransactionDTO{
		Entries: utils.OptionalNullable[[]usecases.UpdateEntryDTO]{Set: true, Value: &updatedEntries},
	})
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to sync transaction entries")
	}

	return nil
}
