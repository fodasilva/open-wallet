package transactions

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/constants"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	transactionRepo "github.com/felipe1496/open-wallet/internal/resources/transactions/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
)

type TransactionsUseCase interface {
	ListViewEntries(ctx context.Context, filter *utils.QueryOptsBuilder) ([]transactionRepo.ViewEntry, error)
	CountViewEntries(ctx context.Context, filter *utils.QueryOptsBuilder) (int, error)
	DeleteTransactionById(id string) error
	CreateTransaction(payload CreateTransactionDTO) (transactionRepo.Transaction, error)
	UpdateTransaction(transactionID string, userID string, payload UpdateTransactionDTO) (transactionRepo.Transaction, error)
}

type transactionsUseCaseImpl struct {
	transactionsRepo  transactionRepo.TransactionsRepo
	entriesRepo       transactionRepo.EntriesRepo
	categoriesUseCase categories.CategoriesUseCase
	db                *sql.DB
}

func NewTransactionsUseCase(transactionsRepo transactionRepo.TransactionsRepo, entriesRepo transactionRepo.EntriesRepo, categoriesUseCase categories.CategoriesUseCase, db *sql.DB) TransactionsUseCase {
	return &transactionsUseCaseImpl{
		transactionsRepo:  transactionsRepo,
		entriesRepo:       entriesRepo,
		categoriesUseCase: categoriesUseCase,
		db:                db,
	}
}

func (uc *transactionsUseCaseImpl) ListViewEntries(ctx context.Context, filter *utils.QueryOptsBuilder) ([]transactionRepo.ViewEntry, error) {
	tracer := otel.Tracer("usecase")
	ctx, span := tracer.Start(ctx, "TransactionsUseCase.ListViewEntries")
	defer span.End()

	entries, err := uc.entriesRepo.Select(uc.db, filter)

	if err != nil {
		span.RecordError(err)
		return []transactionRepo.ViewEntry{}, ErrFailedToFetchEntries
	}

	return entries, nil
}

func (uc *transactionsUseCaseImpl) CountViewEntries(ctx context.Context, filter *utils.QueryOptsBuilder) (int, error) {
	tracer := otel.Tracer("usecase")
	ctx, span := tracer.Start(ctx, "TransactionsUseCase.CountViewEntries")
	defer span.End()

	count, err := uc.entriesRepo.Count(uc.db, filter)

	if err != nil {
		span.RecordError(err)
		return 0, ErrToCountEntries
	}

	return count, nil
}

func (uc *transactionsUseCaseImpl) DeleteTransactionById(id string) error {
	transactionExists, err := uc.transactionsRepo.Select(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return AnErrorOccuredWhileFetchingTransactions
	}

	if len(transactionExists) == 0 {
		return TransactionNotFound
	}

	err = uc.transactionsRepo.Delete(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return ItWasNotPossibleDeleteTransactionErr
	}

	return nil
}

type validateTransactionPropsEntry struct {
	Amount        float64
	ReferenceDate string
}

func validateTransaction(entries []validateTransactionPropsEntry, transactionType constants.TransactionType) error {
	switch transactionType {
	case constants.SimpleExpense:
		{
			if len(entries) > 1 {
				return utils.NewHTTPError(http.StatusBadRequest, "expense must have only one entry")
			}
		}
	case constants.Income:
		{
			if len(entries) > 1 {
				return utils.NewHTTPError(http.StatusBadRequest, "income must have only one entry")
			}
		}
	case constants.Installment:
		{
			if len(entries) < 2 {
				return utils.NewHTTPError(http.StatusBadRequest, "installment must have at least two entries")
			}
		}
	}

	for i, refEntry := range entries {
		iRefDate, _ := time.Parse("2006-01-02", refEntry.ReferenceDate)
		iPeriod := iRefDate.Format("200601")
		for j, currEntry := range entries {
			if i != j {
				jRefDate, _ := time.Parse("2006-01-02", currEntry.ReferenceDate)
				jPeriod := jRefDate.Format("200601")
				if iPeriod == jPeriod {

					return utils.NewHTTPError(http.StatusBadRequest, "entries must be in different periods")
				}
			}
		}

		switch transactionType {
		case constants.Installment:
			{
				if refEntry.Amount >= 0 {
					return utils.NewHTTPError(http.StatusBadRequest, "installment entries must have amount lower than zero")
				}
			}
		case constants.SimpleExpense:
			{
				if refEntry.Amount >= 0 {
					return utils.NewHTTPError(http.StatusBadRequest, "expense entry must have amount lower than zero")
				}
			}
		case constants.Income:
			{
				if refEntry.Amount <= 0 {
					return utils.NewHTTPError(http.StatusBadRequest, "income entry must have amount greater than zero")
				}
			}
		case constants.Recurrence:
			{
				if refEntry.Amount >= 0 {
					return utils.NewHTTPError(http.StatusBadRequest, "recurrence entries must have amount lower than zero")
				}
			}
		}
	}
	return nil
}

func (uc *transactionsUseCaseImpl) CreateTransaction(payload CreateTransactionDTO) (transactionRepo.Transaction, error) {
	err := validateTransaction(func() []validateTransactionPropsEntry {
		entries := make([]validateTransactionPropsEntry, 0)
		if payload.Entries != nil {
			for _, entry := range payload.Entries {
				entries = append(entries, validateTransactionPropsEntry{
					Amount:        entry.Amount,
					ReferenceDate: entry.ReferenceDate,
				})
			}
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
		Note:         utils.OptionalNullable[string]{Set: payload.Note != nil, Value: payload.Note},
		CategoryID:   utils.OptionalNullable[string]{Set: payload.CategoryID != nil, Value: payload.CategoryID},
		RecurrenceID: utils.OptionalNullable[string]{Set: payload.RecurrenceID != nil, Value: payload.RecurrenceID},
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

func (uc *transactionsUseCaseImpl) UpdateTransaction(transactionID string, userID string, payload UpdateTransactionDTO) (t transactionRepo.Transaction, err error) {
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

	if payload.CategoryID != nil && utils.Contains(payload.Update, "category_id") {
		categoryExists, err := uc.categoriesUseCase.List(utils.QueryOpts().
			And("id", "eq", *payload.CategoryID).
			And("user_id", "eq", userID))
		if err != nil {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}

	if utils.ContainsSome(payload.Update, []string{"name", "note", "category_id", "recurrence_id"}) {
		err = uc.transactionsRepo.Update(tx, transactionRepo.UpdateTransactionDTO{
			Name:         utils.OptionalNullable[string]{Set: utils.Contains(payload.Update, "name"), Value: payload.Name},
			Note:         utils.OptionalNullable[string]{Set: utils.Contains(payload.Update, "note"), Value: payload.Note},
			CategoryID:   utils.OptionalNullable[string]{Set: utils.Contains(payload.Update, "category_id"), Value: payload.CategoryID},
			RecurrenceID: utils.OptionalNullable[string]{Set: utils.Contains(payload.Update, "recurrence_id"), Value: payload.RecurrenceID},
		}, utils.QueryOpts().And("id", "eq", transactionID))
		if err != nil {
			return transactionRepo.Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update transaction")
		}
	}

	if payload.Entries != nil && utils.Contains(payload.Update, "entries") {
		err = validateTransaction(func() []validateTransactionPropsEntry {
			entries := make([]validateTransactionPropsEntry, 0)
			if payload.Entries != nil {
				for _, entry := range *payload.Entries {
					entries = append(entries, validateTransactionPropsEntry{
						Amount:        entry.Amount,
						ReferenceDate: entry.ReferenceDate,
					})
				}
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

		for _, entry := range *payload.Entries {
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
