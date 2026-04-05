package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *RecurrencesUseCasesImpl) Update(id string, userID string, payload repository.UpdateRecurrenceDTO) (repository.Recurrence, error) {
	exists, err := uc.repo.Select(uc.db, querybuilder.New().
		And("id", "eq", id).
		And("user_id", "eq", userID))

	if err != nil {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if recurrence exists")
	}

	if len(exists) == 0 {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusNotFound, "recurrence not found")
	}

	if err := uc.validateCategory(userID, payload.CategoryID); err != nil {
		return repository.Recurrence{}, err
	}

	err = uc.repo.Update(uc.db, payload, querybuilder.New().And("id", "eq", id))
	if err != nil {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update recurrence")
	}

	rec, err := uc.fetchRecurrence(id)
	if err != nil {
		return repository.Recurrence{}, err
	}

	if payload.Amount.Set && payload.Amount.Value != nil {
		if err := uc.syncLinkedTransactions(id, userID, *payload.Amount.Value); err != nil {
			return rec, err
		}
	}

	return rec, nil
}

func (uc *RecurrencesUseCasesImpl) validateCategory(userID string, categoryID utils.OptionalNullable[string]) error {
	if !categoryID.Set || categoryID.Value == nil {
		return nil
	}

	exists, err := uc.categoriesUseCase.List(querybuilder.New().
		And("id", "eq", *categoryID.Value).
		And("user_id", "eq", userID))
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
	}

	if len(exists) == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	return nil
}

func (uc *RecurrencesUseCasesImpl) fetchRecurrence(id string) (repository.Recurrence, error) {
	recs, err := uc.repo.Select(uc.db, querybuilder.New().And("id", "eq", id))
	if err != nil || len(recs) == 0 {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch recurrence")
	}
	return recs[0], nil
}

func (uc *RecurrencesUseCasesImpl) syncLinkedTransactions(id string, userID string, amount float64) error {
	txs, err := uc.transactionsUseCase.ListEntries(context.TODO(), querybuilder.New().
		And("user_id", "eq", userID).
		And("recurrence_id", "eq", id))
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

	_, err = uc.transactionsUseCase.UpdateTransaction(transactionID, userID, usecases.UpdateTransactionDTO{
		Entries: utils.OptionalNullable[[]usecases.UpdateEntryDTO]{Set: true, Value: &updatedEntries},
	})
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to sync transaction entries")
	}

	return nil
}
