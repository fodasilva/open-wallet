package usecases

import (
	"database/sql"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

type CategoriesUseCases interface {
	Create(payload repository.CreateCategoryDTO) (repository.Category, error)
	List(filter *querybuilder.Builder) ([]repository.Category, error)
	DeleteByID(id string, userID string) error
	Count(filter *querybuilder.Builder) (int, error)
	ListCategoryAmountPerPeriod(period string, filter *querybuilder.Builder) ([]repository.CategoryAmountPerPeriod, error)
	CountCategoryAmountPerPeriod(period string, filter *querybuilder.Builder) (int, error)
	Update(id string, userID string, payload repository.UpdateCategoryDTO) (repository.Category, error)
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
