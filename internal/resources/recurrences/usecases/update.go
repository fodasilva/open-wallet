package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *RecurrencesUseCasesImpl) Update(id string, userID string, payload repository.UpdateRecurrenceDTO) (repository.Recurrence, error) {
	exists, err := uc.repo.Select(uc.db, utils.QueryOpts().
		And("id", "eq", id).
		And("user_id", "eq", userID))

	if err != nil {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if recurrence exists")
	}

	if len(exists) == 0 {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusNotFound, "recurrence not found")
	}

	if payload.CategoryID.Set && payload.CategoryID.Value != nil {
		categoryExists, err := uc.categoriesUseCase.List(utils.QueryOpts().
			And("id", "eq", *payload.CategoryID.Value).
			And("user_id", "eq", userID))
		if err != nil {
			return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return repository.Recurrence{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}

	err = uc.repo.Update(uc.db, payload, utils.QueryOpts().And("id", "eq", id))
	if err != nil {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update recurrence")
	}

	updatedRecs, err := uc.repo.Select(uc.db, utils.QueryOpts().And("id", "eq", id))
	if err != nil || len(updatedRecs) == 0 {
		return repository.Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch updated recurrence")
	}
	rec := updatedRecs[0]

	if payload.Amount.Set && payload.Amount.Value != nil {
		txs, err := uc.transactionsUseCase.ListViewEntries(context.TODO(), utils.QueryOpts().
			And("user_id", "eq", userID).
			And("recurrence_id", "eq", id))
		if err != nil {
			return rec, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch linked transactions for sync")
		}

		if len(txs) > 0 {
			transactionID := txs[0].TransactionID
			var updatedEntries []transactions.UpdateEntryDTO
			for _, entry := range txs {
				updatedEntries = append(updatedEntries, transactions.UpdateEntryDTO{
					Amount:        *payload.Amount.Value,
					ReferenceDate: entry.ReferenceDate,
				})
			}

			_, err = uc.transactionsUseCase.UpdateTransaction(transactionID, userID, transactions.UpdateTransactionDTO{
				Update:  []string{"entries"},
				Entries: &updatedEntries,
			})
			if err != nil {
				return rec, utils.NewHTTPError(http.StatusInternalServerError, "failed to sync transaction entries")
			}
		}
	}

	return rec, nil
}
