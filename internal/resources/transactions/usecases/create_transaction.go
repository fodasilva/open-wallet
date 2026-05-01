package usecases

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/oklog/ulid/v2"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (uc *TransactionsUseCasesImpl) CreateTransaction(ctx context.Context, payload CreateTransactionDTO) (transactionRepo.Transaction, error) {
	if err := uc.validateCategory(ctx, payload.UserID, payload.CategoryID); err != nil {
		return transactionRepo.Transaction{}, err
	}

	if err := uc.validatePayloadEntries(payload.Entries, payload.Type); err != nil {
		return transactionRepo.Transaction{}, err
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return transactionRepo.Transaction{}, httputil.NewHTTPError(http.StatusInternalServerError, "failed to begin transaction")
	}

	transactionID := ulid.Make().String()
	err = uc.transactionsRepo.Insert(ctx, tx, transactionRepo.CreateTransactionDTO{
		ID:           transactionID,
		UserID:       payload.UserID,
		Type:         payload.Type,
		Name:         payload.Name,
		Note:         payload.Note,
		CategoryID:   payload.CategoryID,
		RecurrenceID: payload.RecurrenceID,
	})

	if err != nil {
		_ = tx.Rollback()
		return transactionRepo.Transaction{}, httputil.NewHTTPError(http.StatusInternalServerError, "failed to create transaction")
	}

	if err := uc.persistEntries(ctx, tx, transactionID, payload.Entries, payload.Type); err != nil {
		_ = tx.Rollback()
		return transactionRepo.Transaction{}, err
	}

	if err := tx.Commit(); err != nil {
		return transactionRepo.Transaction{}, httputil.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}

	return uc.fetchCreatedTransaction(ctx, transactionID)
}

func (uc *TransactionsUseCasesImpl) validateCategory(ctx context.Context, userID string, categoryID util.OptionalNullable[string]) error {
	if !categoryID.Set || categoryID.Value == nil {
		return nil
	}

	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().
		And("id", "eq", *categoryID.Value).
		And("user_id", "eq", userID))
	exists, err := uc.categoriesUseCase.List(filterCtx)
	if err != nil {
		return httputil.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
	}

	if len(exists) == 0 {
		return httputil.NewHTTPError(http.StatusNotFound, "category not found")
	}

	return nil
}

func (uc *TransactionsUseCasesImpl) validatePayloadEntries(entries []CreateEntryDTO, tType transactionRepo.TransactionType) error {
	props := make([]validateTransactionPropsEntry, len(entries))
	for i, entry := range entries {
		props[i] = validateTransactionPropsEntry(entry)
	}
	return validateTransaction(props, tType)
}

func (uc *TransactionsUseCasesImpl) persistEntries(ctx context.Context, tx *sql.Tx, transactionID string, entries []CreateEntryDTO, tType transactionRepo.TransactionType) error {
	for _, entry := range entries {
		amount := entry.Amount
		if (tType == transactionRepo.SimpleExpense || tType == transactionRepo.Installment || tType == transactionRepo.Recurrence) && amount > 0 {
			amount *= -1
		} else if tType == transactionRepo.Income && amount < 0 {
			amount *= -1
		}

		err := uc.entriesRepo.Insert(ctx, tx, transactionRepo.CreateEntryDTO{
			ID:            ulid.Make().String(),
			TransactionID: transactionID,
			Amount:        amount,
			ReferenceDate: entry.ReferenceDate,
		})

		if err != nil {
			return httputil.NewHTTPError(http.StatusInternalServerError, "failed to create entry")
		}
	}
	return nil
}

func (uc *TransactionsUseCasesImpl) fetchCreatedTransaction(ctx context.Context, id string) (transactionRepo.Transaction, error) {
	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", id))
	created, err := uc.transactionsRepo.Select(filterCtx, uc.db)
	if err != nil || len(created) == 0 {
		return transactionRepo.Transaction{}, httputil.NewHTTPError(http.StatusInternalServerError, "failed to fetch created transaction")
	}
	return created[0], nil
}
