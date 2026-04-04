package usecases

import (
	"net/http"

	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/oklog/ulid/v2"
)

func (uc *TransactionsUseCasesImpl) UpdateTransaction(transactionID string, userID string, payload UpdateTransactionDTO) (t transactionRepo.Transaction, err error) {
	tx, err := uc.db.Begin()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if err != nil {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to start transaction")
	}

	exists, err := uc.entriesRepo.Select(tx, utils.QueryOpts().
		And("transaction_id", "eq", transactionID).
		And("user_id", "eq", userID))
	if err != nil {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if transaction exists")
	}

	if len(exists) == 0 {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusNotFound, "transaction not found")
	}

	if payload.CategoryID.Set && payload.CategoryID.Value != nil {
		categoryExists, err := uc.categoriesUseCase.List(utils.QueryOpts().
			And("id", "eq", *payload.CategoryID.Value).
			And("user_id", "eq", userID))
		if err != nil {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}

	if payload.Name.Set || payload.Note.Set || payload.CategoryID.Set || payload.RecurrenceID.Set {
		err = uc.transactionsRepo.Update(tx, transactionRepo.UpdateTransactionDTO{
			Name:         payload.Name,
			Note:         payload.Note,
			CategoryID:   payload.CategoryID,
			RecurrenceID: payload.RecurrenceID,
		}, utils.QueryOpts().And("id", "eq", transactionID))
		if err != nil {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update transaction")
		}
	}

	if payload.Entries.Set && payload.Entries.Value != nil {
		err = validateTransaction(func() []validateTransactionPropsEntry {
			entries := make([]validateTransactionPropsEntry, 0)
			for _, entry := range *payload.Entries.Value {
				entries = append(entries, validateTransactionPropsEntry{
					Amount:        entry.Amount,
					ReferenceDate: entry.ReferenceDate,
				})
			}
			return entries
		}(), exists[0].Type)

		if err != nil {
			return transactionRepo.Transaction{}, err
		}

		err = uc.entriesRepo.Delete(tx, utils.QueryOpts().And("transaction_id", "eq", exists[0].TransactionID))
		if err != nil {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to delete previous entries")
		}

		for _, entry := range *payload.Entries.Value {
			err = uc.entriesRepo.Insert(tx, transactionRepo.CreateEntryDTO{
				ID:            ulid.Make().String(),
				TransactionID: exists[0].TransactionID,
				Amount:        entry.Amount,
				ReferenceDate: entry.ReferenceDate,
			})
			if err != nil {
				return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create entry")
			}
		}
	}

	transactions, err := uc.transactionsRepo.Select(tx, utils.QueryOpts().And("id", "eq", transactionID))
	if err != nil || len(transactions) == 0 {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to list transaction")
	}

	return transactions[0], nil
}
