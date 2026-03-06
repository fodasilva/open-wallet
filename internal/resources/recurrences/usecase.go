package recurrences

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/constants"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"
	"go.opentelemetry.io/otel"
)

type RecurrencesUseCase interface {
	Create(payload CreateRecurrenceDTO) (Recurrence, error)
	List(ctx context.Context, filter *utils.QueryOptsBuilder) ([]Recurrence, error)
	Count(ctx context.Context, filter *utils.QueryOptsBuilder) (int, error)
	DeleteByID(id string, scope string) error
	Update(id string, userID string, payload UpdateRecurrenceDTO) (Recurrence, error)
	PrepareRecurrences(ctx context.Context, userID string, targetPeriod string) error
}

type recurrencesUseCaseImpl struct {
	repo                RecurrencesRepo
	categoriesUseCase   categories.CategoriesUseCase
	transactionsUseCase transactions.TransactionsUseCase
	db                  *sql.DB
}

func NewRecurrencesUseCase(repo RecurrencesRepo, categoriesUseCase categories.CategoriesUseCase, transactionsUseCase transactions.TransactionsUseCase, db *sql.DB) RecurrencesUseCase {
	return &recurrencesUseCaseImpl{
		repo:                repo,
		categoriesUseCase:   categoriesUseCase,
		transactionsUseCase: transactionsUseCase,
		db:                  db,
	}
}

func (uc *recurrencesUseCaseImpl) Create(payload CreateRecurrenceDTO) (Recurrence, error) {
	if payload.CategoryID != nil {
		categoryExists, err := uc.categoriesUseCase.List(utils.QueryOpts().
			And("id", "eq", *payload.CategoryID).
			And("user_id", "eq", payload.UserID))
		if err != nil {
			return Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return Recurrence{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}

	rec, err := uc.repo.Create(uc.db, payload)
	if err != nil {
		return Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create recurrence")
	}

	return rec, nil
}

func (uc *recurrencesUseCaseImpl) List(ctx context.Context, filter *utils.QueryOptsBuilder) ([]Recurrence, error) {
	tracer := otel.Tracer("usecase")
	ctx, span := tracer.Start(ctx, "RecurrencesUseCase.List")
	defer span.End()

	items, err := uc.repo.List(ctx, uc.db, filter)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to list recurrences")
	}

	return items, nil
}

func (uc *recurrencesUseCaseImpl) Count(ctx context.Context, filter *utils.QueryOptsBuilder) (int, error) {
	tracer := otel.Tracer("usecase")
	ctx, span := tracer.Start(ctx, "RecurrencesUseCase.Count")
	defer span.End()

	count, err := uc.repo.Count(ctx, uc.db, filter)
	if err != nil {
		span.RecordError(err)
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count recurrences")
	}

	return count, nil
}

func (uc *recurrencesUseCaseImpl) DeleteByID(id string, scope string) error {
	exists, err := uc.repo.List(context.TODO(), uc.db, utils.QueryOpts().And("id", "eq", id))
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

	err = uc.repo.DeleteByID(uc.db, id)
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete recurrence")
	}

	return nil
}

func (uc *recurrencesUseCaseImpl) Update(id string, userID string, payload UpdateRecurrenceDTO) (Recurrence, error) {
	exists, err := uc.repo.List(context.TODO(), uc.db, utils.QueryOpts().
		And("id", "eq", id).
		And("user_id", "eq", userID))

	if err != nil {
		return Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if recurrence exists")
	}

	if len(exists) == 0 {
		return Recurrence{}, utils.NewHTTPError(http.StatusNotFound, "recurrence not found")
	}

	if payload.CategoryID != nil && utils.Contains(payload.Update, "category_id") {
		categoryExists, err := uc.categoriesUseCase.List(utils.QueryOpts().
			And("id", "eq", *payload.CategoryID).
			And("user_id", "eq", userID))
		if err != nil {
			return Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return Recurrence{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}

	rec, err := uc.repo.Update(uc.db, id, payload)
	if err != nil {
		return Recurrence{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update recurrence")
	}

	if payload.Amount != nil && utils.Contains(payload.Update, "amount") {
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
					Amount:        *payload.Amount,
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

func (uc *recurrencesUseCaseImpl) PrepareRecurrences(ctx context.Context, userID string, targetPeriod string) error {
	recurrences, err := uc.repo.List(ctx, uc.db, utils.QueryOpts().And("user_id", "eq", userID))
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
