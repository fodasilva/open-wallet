package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/constants"
	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/oklog/ulid/v2"
)

func (uc *TransactionsUseCasesImpl) CreateTransaction(payload CreateTransactionDTO) (transactionRepo.Transaction, error) {
	if payload.CategoryID.Set && payload.CategoryID.Value != nil {
		categoryExists, err := uc.categoriesUseCase.List(utils.QueryOpts().
			And("id", "eq", *payload.CategoryID.Value).
			And("user_id", "eq", payload.UserID))
		if err != nil {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}
	err := validateTransaction(func() []validateTransactionPropsEntry {
		entries := make([]validateTransactionPropsEntry, 0)
		for _, entry := range payload.Entries {
			entries = append(entries, validateTransactionPropsEntry{
				Amount:        entry.Amount,
				ReferenceDate: entry.ReferenceDate,
			})
		}
		return entries
	}(), payload.Type)
	if err != nil {
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
		tx.Rollback()
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create transaction")
	}

	for _, entry := range payload.Entries {
		if (payload.Type == constants.SimpleExpense || payload.Type == constants.Installment) && entry.Amount > 0 {
			entry.Amount = entry.Amount * -1
		} else if payload.Type == constants.Income && entry.Amount < 0 {
			entry.Amount = entry.Amount * -1
		}

		err = uc.entriesRepo.Insert(tx, transactionRepo.CreateEntryDTO{
			ID:            ulid.Make().String(),
			TransactionID: transactionID,
			Amount:        entry.Amount,
			ReferenceDate: entry.ReferenceDate,
		})

		if err != nil {
			tx.Rollback()
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create entry")
		}
	}

	err = tx.Commit()

	if err != nil {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}

	// Always fetch after mutation
	created, err := uc.transactionsRepo.Select(uc.db, utils.QueryOpts().And("id", "eq", transactionID))
	if err != nil || len(created) == 0 {
		return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch created transaction")
	}

	return created[0], nil
}
