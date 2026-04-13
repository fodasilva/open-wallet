package usecases

import (
	"context"
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *RecurrencesUseCasesImpl) DeleteByID(ctx context.Context, id string, userID string, scope string) error {
	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", id).And("user_id", "eq", userID))
	exists, err := uc.repo.Select(filterCtx, uc.db)
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch recurrence")
	}

	if len(exists) == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "recurrence not found")
	}

	rec := exists[0]

	// Find linked transaction
	txFilterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().
		And("user_id", "eq", rec.UserID).
		And("recurrence_id", "eq", rec.ID))
	txs, err := uc.transactionsUseCase.ListEntries(txFilterCtx)

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch linked transactions")
	}

	if len(txs) > 0 {
		transactionID := txs[0].TransactionID
		switch scope {
		case "all":
			err = uc.transactionsUseCase.DeleteTransactionById(ctx, transactionID, rec.UserID)
			if err != nil {
				return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete linked transaction")
			}
		case "until_current":
			currentPeriod := time.Now().Format("200601")
			var filteredEntries []usecases.UpdateEntryDTO
			for _, entry := range txs {
				if entry.Period <= currentPeriod {
					filteredEntries = append(filteredEntries, usecases.UpdateEntryDTO{
						Amount:        entry.Amount,
						ReferenceDate: entry.ReferenceDate,
					})
				}
			}

			_, err = uc.transactionsUseCase.UpdateTransaction(ctx, transactionID, rec.UserID, usecases.UpdateTransactionDTO{
				Entries:      utils.OptionalNullable[[]usecases.UpdateEntryDTO]{Set: true, Value: &filteredEntries},
				RecurrenceID: utils.OptionalNullable[string]{Set: true, Value: nil},
			})
			if err != nil {
				return utils.NewHTTPError(http.StatusInternalServerError, "failed to update transactions entries")
			}
		}
	}

	err = uc.repo.Delete(filterCtx, uc.db)
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete recurrence")
	}

	return nil
}
