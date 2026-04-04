package usecases

import (
	"context"
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *RecurrencesUseCasesImpl) DeleteByID(id string, scope string) error {
	exists, err := uc.repo.Select(uc.db, utils.QueryOpts().And("id", "eq", id))
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch recurrence")
	}

	if len(exists) == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "recurrence not found")
	}

	rec := exists[0]

	// Find linked transaction
	txs, err := uc.transactionsUseCase.ListViewEntries(context.TODO(), utils.QueryOpts().
		And("user_id", "eq", rec.UserID).
		And("recurrence_id", "eq", rec.ID))

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch linked transactions")
	}

	if len(txs) > 0 {
		transactionID := txs[0].TransactionID
		switch scope {
		case "all":
			err = uc.transactionsUseCase.DeleteTransactionById(transactionID)
			if err != nil {
				return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete linked transaction")
			}
		case "until_current":
			currentPeriod := time.Now().Format("200601")
			var filteredEntries []transactions.UpdateEntryDTO
			for _, entry := range txs {
				if entry.Period <= currentPeriod {
					filteredEntries = append(filteredEntries, transactions.UpdateEntryDTO{
						Amount:        entry.Amount,
						ReferenceDate: entry.ReferenceDate,
					})
				}
			}

			_, err = uc.transactionsUseCase.UpdateTransaction(transactionID, rec.UserID, transactions.UpdateTransactionDTO{
				Update:       []string{"entries", "recurrence_id"},
				Entries:      &filteredEntries,
				RecurrenceID: nil,
			})
			if err != nil {
				return utils.NewHTTPError(http.StatusInternalServerError, "failed to update transactions entries")
			}
		}
	}

	err = uc.repo.Delete(uc.db, utils.QueryOpts().And("id", "eq", id))
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete recurrence")
	}

	return nil
}
