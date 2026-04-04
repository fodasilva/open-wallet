package usecases

import (
	"context"
	"fmt"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"net/http"
	"time"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *RecurrencesUseCasesImpl) PrepareRecurrences(ctx context.Context, userID string, targetPeriod string) error {
	recurrences, err := uc.repo.Select(uc.db, utils.QueryOpts().And("user_id", "eq", userID))
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch recurrences")
	}

	for _, rec := range recurrences {
		existingTxs, err := uc.transactionsUseCase.ListEntries(ctx, utils.QueryOpts().
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
			_, err := uc.transactionsUseCase.CreateTransaction(usecases.CreateTransactionDTO{
				UserID:       userID,
				Name:         rec.Name,
				CategoryID:   utils.OptionalNullable[string]{Set: rec.CategoryID != nil, Value: rec.CategoryID},
				Note:         utils.OptionalNullable[string]{Set: rec.Note != nil, Value: rec.Note},
				Type:         transactionRepo.Recurrence,
				RecurrenceID: utils.OptionalNullable[string]{Set: true, Value: &rec.ID},
				Entries: []usecases.CreateEntryDTO{{
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

			var payloadEntries []usecases.UpdateEntryDTO
			for _, t := range existingTxs {
				payloadEntries = append(payloadEntries, usecases.UpdateEntryDTO{
					Amount:        t.Amount,
					ReferenceDate: t.ReferenceDate,
				})
			}

			payloadEntries = append(payloadEntries, usecases.UpdateEntryDTO{
				Amount:        rec.Amount,
				ReferenceDate: newDateStr,
			})

			_, err = uc.transactionsUseCase.UpdateTransaction(targetTxID, userID, usecases.UpdateTransactionDTO{
				Entries: utils.OptionalNullable[[]usecases.UpdateEntryDTO]{Set: true, Value: &payloadEntries},
			})

			if err != nil {
				return err
			}
		}
	}
	return nil
}
