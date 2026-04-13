package usecases

import (
	"context"
	"database/sql"

	"github.com/felipe1496/open-wallet/internal/resources/categories/usecases"
	"github.com/felipe1496/open-wallet/internal/resources/recurrences/repository"
	transactionsUseCases "github.com/felipe1496/open-wallet/internal/resources/transactions/usecases"
)

type RecurrencesUseCases interface {
	Create(ctx context.Context, payload repository.CreateRecurrenceDTO) (repository.Recurrence, error)
	List(ctx context.Context) ([]repository.Recurrence, error)
	Count(ctx context.Context) (int, error)
	DeleteByID(ctx context.Context, id string, userID string, scope string) error
	Update(ctx context.Context, id string, userID string, payload repository.UpdateRecurrenceDTO) (repository.Recurrence, error)
	PrepareRecurrences(ctx context.Context, userID string, targetPeriod string) error
}

type RecurrencesUseCasesImpl struct {
	repo                repository.RecurrencesRepo
	categoriesUseCase   usecases.CategoriesUseCases
	transactionsUseCase transactionsUseCases.TransactionsUseCases
	db                  *sql.DB
}

func NewRecurrencesUseCases(
	repo repository.RecurrencesRepo,
	categoriesUseCase usecases.CategoriesUseCases,
	transactionsUseCase transactionsUseCases.TransactionsUseCases,
	db *sql.DB,
) RecurrencesUseCases {

	return &RecurrencesUseCasesImpl{
		repo:                repo,
		categoriesUseCase:   categoriesUseCase,
		transactionsUseCase: transactionsUseCase,
		db:                  db,
	}
}
