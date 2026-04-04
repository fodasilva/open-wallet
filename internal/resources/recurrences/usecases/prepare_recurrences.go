package usecases

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/constants"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *RecurrencesUseCasesImpl) PrepareRecurrences(ctx context.Context, userID string, targetPeriod string) error {
	recurrences, err := uc.repo.Select(uc.db, utils.QueryOpts().And("user_id", "eq", userID))
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch recurrences")
	}

	for _, rec := range recurrences {
		existingTxs, err := uc.transactionsUseCase.ListViewEntries(ctx, utils.QueryOpts().
			And("user_id", "eq", userID).
			And("recurrence_id", "eq", rec.ID))
		if err != nil {
			return utils.NewHTTPError(http.StatusInternalServerError, "failed to check existing transactions")
		}

		var targetTxID string

		newDateStr := fmt.Sprintf("%s-%s-%02d", targetPeriod[:4], targetPeriod[4:], rec.DayOfMonth)
		if _, err := time.Parse("2006-01-02", newDateStr); err != nil {
			newDateStr = fmt.Sprintf("%s-%s-01", targetPeriod[:4], targetPeriod[4:])
		}

		if len(existingTxs) == 0 {
			_, err := uc.transactionsUseCase.CreateTransaction(transactions.CreateTransactionDTO{
				UserID:       userID,
				Name:         rec.Name,
				CategoryID:   rec.CategoryID,
				Note:         rec.Note,
				Type:         constants.Recurrence,
				RecurrenceID: &rec.ID,
				Entries: []transactions.CreateEntryDTO{{
					Amount:        rec.Amount,
					ReferenceDate: newDateStr,
				}},
			})
			if err != nil {
				return err
			}
		} else {
			targetTxID = existingTxs[0].TransactionID

			hasPeriod := false
			for _, t := range existingTxs {
				if t.Period == targetPeriod {
					hasPeriod = true
					break
				}
			}

			if hasPeriod {
				continue
			}

			var payloadEntries []transactions.UpdateEntryDTO
			for _, t := range existingTxs {
				payloadEntries = append(payloadEntries, transactions.UpdateEntryDTO{
					Amount:        t.Amount,
					ReferenceDate: t.ReferenceDate,
				})
			}

			payloadEntries = append(payloadEntries, transactions.UpdateEntryDTO{
				Amount:        rec.Amount,
				ReferenceDate: newDateStr,
			})

			_, err = uc.transactionsUseCase.UpdateTransaction(targetTxID, userID, transactions.UpdateTransactionDTO{
				Update:  []string{"entries"},
				Entries: &payloadEntries,
			})

			if err != nil {
				return err
			}
		}
	}
	return nil
}
