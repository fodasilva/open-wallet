package usecases

import (
	"database/sql"
	"net/http"

	"github.com/oklog/ulid/v2"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *TransactionsUseCasesImpl) UpdateTransaction(transactionID string, userID string, payload UpdateTransactionDTO) (t transactionRepo.Transaction, err error) {
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

	exists, err := uc.entriesRepo.Select(tx, utils.QueryOpts().
		And("transaction_id", "eq", transactionID).
		And("user_id", "eq", userID))
	if err != nil {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if transaction exists")
	}

	if len(exists) == 0 {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusNotFound, "transaction not found")
	}

	if err := uc.validateCategory(userID, payload.CategoryID); err != nil {
		return transactionRepo.Transaction{}, err
	}

	if err := uc.updateTransactionMetadata(tx, transactionID, payload); err != nil {
		return transactionRepo.Transaction{}, err
	}

	if err := uc.syncTransactionEntries(tx, transactionID, exists[0].Type, payload.Entries); err != nil {
		return transactionRepo.Transaction{}, err
	}

	created, err := uc.transactionsRepo.Select(tx, utils.QueryOpts().And("id", "eq", transactionID))
	if err != nil || len(created) == 0 {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch updated transaction")
	}

	return created[0], nil
}

func (uc *TransactionsUseCasesImpl) updateTransactionMetadata(tx *sql.Tx, transactionID string, payload UpdateTransactionDTO) error {
	if !payload.Name.Set && !payload.Note.Set && !payload.CategoryID.Set && !payload.RecurrenceID.Set {
		return nil
	}

	err := uc.transactionsRepo.Update(tx, transactionRepo.UpdateTransactionDTO{
		Name:         payload.Name,
		Note:         payload.Note,
		CategoryID:   payload.CategoryID,
		RecurrenceID: payload.RecurrenceID,
	}, utils.QueryOpts().And("id", "eq", transactionID))

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to update transaction")
	}

	return nil
}

func (uc *TransactionsUseCasesImpl) syncTransactionEntries(tx *sql.Tx, transactionID string, tType transactionRepo.TransactionType, entriesOpt utils.OptionalNullable[[]UpdateEntryDTO]) error {
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

	err := uc.entriesRepo.Delete(tx, utils.QueryOpts().And("transaction_id", "eq", transactionID))
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete previous entries")
	}

	for _, entry := range entries {
		err = uc.entriesRepo.Insert(tx, transactionRepo.CreateEntryDTO{
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
