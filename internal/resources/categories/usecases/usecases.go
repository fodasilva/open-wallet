package usecases

import (
	"database/sql"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type CategoriesUseCases interface {
	Create(payload repository.CreateCategoryDTO) (repository.Category, error)
	List(filter *utils.QueryOptsBuilder) ([]repository.Category, error)
	DeleteByID(id string) error
	Count(filter *utils.QueryOptsBuilder) (int, error)
	ListCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) ([]repository.CategoryAmountPerPeriod, error)
	CountCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) (int, error)
	Update(id string, payload repository.UpdateCategoryDTO) (repository.Category, error)
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
