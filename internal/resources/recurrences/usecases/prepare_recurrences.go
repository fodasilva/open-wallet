package usecases

import (
	"context"
	"fmt"
	"net/http"
	"time"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (uc *RecurrencesUseCasesImpl) PrepareRecurrences(ctx context.Context, userID string, targetPeriod string) error {
	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("user_id", "eq", userID))
	recurrences, err := uc.repo.Select(filterCtx, uc.db)
	if err != nil {
		return util.NewHTTPError(http.StatusInternalServerError, "failed to fetch recurrences")
	}

	for _, rec := range recurrences {
		if rec.StartPeriod > targetPeriod {
			continue
		}
		if rec.EndPeriod != nil && *rec.EndPeriod < targetPeriod {
			continue
		}

		txFilterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().
			And("user_id", "eq", userID).
			And("recurrence_id", "eq", rec.ID))
		existingTxs, err := uc.transactionsUseCase.ListEntries(txFilterCtx)
		if err != nil {
			return util.NewHTTPError(http.StatusInternalServerError, "failed to check existing transactions")
		}

		var targetTxID string

		newDateStr := fmt.Sprintf("%s-%s-%02d", targetPeriod[:4], targetPeriod[4:], rec.DayOfMonth)
		newDate, err := time.Parse("2006-01-02", newDateStr)
		if err != nil {
			newDate, _ = time.Parse("2006-01-02", fmt.Sprintf("%s-%s-01", targetPeriod[:4], targetPeriod[4:]))
		}

		if len(existingTxs) == 0 {
			_, err := uc.transactionsUseCase.CreateTransaction(ctx, usecases.CreateTransactionDTO{
				UserID:       userID,
				Name:         rec.Name,
				CategoryID:   util.OptionalNullable[string]{Set: rec.CategoryID != nil, Value: rec.CategoryID},
				Note:         util.OptionalNullable[string]{Set: rec.Note != nil, Value: rec.Note},
				Type:         transactionRepo.Recurrence,
				RecurrenceID: util.OptionalNullable[string]{Set: true, Value: &rec.ID},
				Entries: []usecases.CreateEntryDTO{{
					Amount:        rec.Amount,
					ReferenceDate: newDate,
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
				ReferenceDate: newDate,
			})

			_, err = uc.transactionsUseCase.UpdateTransaction(ctx, targetTxID, userID, usecases.UpdateTransactionDTO{
				Entries: util.OptionalNullable[[]usecases.UpdateEntryDTO]{Set: true, Value: &payloadEntries},
			})

			if err != nil {
				return err
			}
		}
	}
	return nil
}
