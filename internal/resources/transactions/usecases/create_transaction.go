package usecases

import (
	"database/sql"
	"net/http"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/oklog/ulid/v2"
)

func (uc *TransactionsUseCasesImpl) CreateTransaction(payload CreateTransactionDTO) (transactionRepo.Transaction, error) {
	if err := uc.validateCategory(payload.UserID, payload.CategoryID); err != nil {
		return transactionRepo.Transaction{}, err
	}

	if err := uc.validatePayloadEntries(payload.Entries, payload.Type); err != nil {
		return transactionRepo.Transaction{}, err
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to begin transaction")
	}

	transactionID := ulid.Make().String()
	err = uc.transactionsRepo.Insert(tx, transactionRepo.CreateTransactionDTO{
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
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create transaction")
	}

	if err := uc.persistEntries(tx, transactionID, payload.Entries, payload.Type); err != nil {
		_ = tx.Rollback()
		return transactionRepo.Transaction{}, err
	}

	if err := tx.Commit(); err != nil {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}

	return uc.fetchCreatedTransaction(transactionID)
}

func (uc *TransactionsUseCasesImpl) validateCategory(userID string, categoryID utils.OptionalNullable[string]) error {
	if !categoryID.Set || categoryID.Value == nil {
		return nil
	}

	exists, err := uc.categoriesUseCase.List(utils.QueryOpts().
		And("id", "eq", *categoryID.Value).
		And("user_id", "eq", userID))
	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
	}

	if len(exists) == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "category not found")
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

func (uc *TransactionsUseCasesImpl) persistEntries(tx *sql.Tx, transactionID string, entries []CreateEntryDTO, tType transactionRepo.TransactionType) error {
	for _, entry := range entries {
		amount := entry.Amount
		if (tType == transactionRepo.SimpleExpense || tType == transactionRepo.Installment) && amount > 0 {
			amount *= -1
		} else if tType == transactionRepo.Income && amount < 0 {
			amount *= -1
		}

		err := uc.entriesRepo.Insert(tx, transactionRepo.CreateEntryDTO{
			ID:            ulid.Make().String(),
			TransactionID: transactionID,
			Amount:        amount,
			ReferenceDate: entry.ReferenceDate,
		})

		if err != nil {
			return utils.NewHTTPError(http.StatusInternalServerError, "failed to create entry")
		}
	}
	return nil
}

func (uc *TransactionsUseCasesImpl) fetchCreatedTransaction(id string) (transactionRepo.Transaction, error) {
	created, err := uc.transactionsRepo.Select(uc.db, utils.QueryOpts().And("id", "eq", id))
	if err != nil || len(created) == 0 {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch created transaction")
	}
	return created[0], nil
}
