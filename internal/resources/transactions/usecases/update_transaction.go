package usecases

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/oklog/ulid/v2"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *TransactionsUseCasesImpl) UpdateTransaction(ctx context.Context, transactionID string, userID string, payload UpdateTransactionDTO) (t transactionRepo.Transaction, err error) {
	tx, err := uc.db.Begin()
	if err != nil {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to start transaction")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().
		And("transaction_id", "eq", transactionID).
		And("user_id", "eq", userID))
	exists, err := uc.entriesRepo.Select(filterCtx, tx)
	if err != nil {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if transaction exists")
	}

	if len(exists) == 0 {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusNotFound, "transaction not found")
	}

	if err := uc.validateCategory(ctx, userID, payload.CategoryID); err != nil {
		return transactionRepo.Transaction{}, err
	}

	if err := uc.updateTransactionMetadata(ctx, tx, transactionID, payload); err != nil {
		return transactionRepo.Transaction{}, err
	}

	if err := uc.syncTransactionEntries(ctx, tx, transactionID, exists[0].Type, payload.Entries); err != nil {
		return transactionRepo.Transaction{}, err
	}

	createdCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", transactionID))
	created, err := uc.transactionsRepo.Select(createdCtx, tx)
	if err != nil || len(created) == 0 {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch updated transaction")
	}

	return created[0], nil
}

func (uc *TransactionsUseCasesImpl) updateTransactionMetadata(ctx context.Context, tx *sql.Tx, transactionID string, payload UpdateTransactionDTO) error {
	if !payload.Name.Set && !payload.Note.Set && !payload.CategoryID.Set && !payload.RecurrenceID.Set {
		return nil
	}

	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", transactionID))
	err := uc.transactionsRepo.Update(filterCtx, tx, transactionRepo.UpdateTransactionDTO{
		Name:         payload.Name,
		Note:         payload.Note,
		CategoryID:   payload.CategoryID,
		RecurrenceID: payload.RecurrenceID,
	})

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to update transaction")
	}

	return nil
}

func (uc *TransactionsUseCasesImpl) syncTransactionEntries(ctx context.Context, tx *sql.Tx, transactionID string, tType transactionRepo.TransactionType, entriesOpt utils.OptionalNullable[[]UpdateEntryDTO]) error {
	if !entriesOpt.Set || entriesOpt.Value == nil {
		return nil
	}

	entries := *entriesOpt.Value
	props := make([]validateTransactionPropsEntry, len(entries))
	for i, entry := range entries {
		props[i] = validateTransactionPropsEntry(entry)
	}

	if err := validateTransaction(props, tType); err != nil {
		return err
	}

	deleteCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("transaction_id", "eq", transactionID))
	err := uc.entriesRepo.Delete(deleteCtx, tx)
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete previous entries")
	}

	for _, entry := range entries {
		err = uc.entriesRepo.Insert(ctx, tx, transactionRepo.CreateEntryDTO{
			ID:            ulid.Make().String(),
			TransactionID: transactionID,
			Amount:        entry.Amount,
			ReferenceDate: entry.ReferenceDate,
		})
		if err != nil {
			return utils.NewHTTPError(http.StatusInternalServerError, "failed to create entry")
		}
	}

	return nil
}
