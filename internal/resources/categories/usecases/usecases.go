package usecases

import (
	"context"
	"database/sql"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
)

type CategoriesUseCases interface {
	Create(ctx context.Context, payload repository.CreateCategoryDTO) (repository.Category, error)
	List(ctx context.Context) ([]repository.Category, error)
	DeleteByID(ctx context.Context, id string, userID string) error
	Count(ctx context.Context) (int, error)
	ListCategoryAmountPerPeriod(ctx context.Context, period string) ([]repository.CategoryAmountPerPeriod, error)
	CountCategoryAmountPerPeriod(ctx context.Context, period string) (int, error)
	Update(ctx context.Context, id string, userID string, payload repository.UpdateCategoryDTO) (repository.Category, error)
}

type CategoriesUseCasesImpl struct {
	repo repository.CategoriesRepo
	db   *sql.DB
}

func NewCategoriesUseCases(repo repository.CategoriesRepo, db *sql.DB) CategoriesUseCases {
	return &CategoriesUseCasesImpl{
		repo: repo,
		db:   db,
	}
}
