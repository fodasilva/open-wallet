package categories

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
)

type CategoriesUseCase interface {
	Create(payload CreateCategoryDTO) (Category, error)
	List(filter *utils.QueryOptsBuilder) ([]Category, error)
	DeleteByID(id string) error
	Count(filter *utils.QueryOptsBuilder) (int, error)
	ListCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) ([]CategoryAmountPerPeriod, error)
	CountCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) (int, error)
	Update(id string, payload UpdateCategoryDTO) (Category, error)
}

type CategoriesUseCaseImpl struct {
	repo CategoriesRepo
	db   *sql.DB
}

func NewCategoriesUseCase(repo CategoriesRepo, db *sql.DB) CategoriesUseCase {
	return &CategoriesUseCaseImpl{
		repo: repo,
		db:   db,
	}
}

func (uc *CategoriesUseCaseImpl) Create(payload CreateCategoryDTO) (Category, error) {
	category, err := uc.repo.Create(uc.db, payload)

	if err != nil {
		return Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create category")
	}

	return category, nil
}

func (uc *CategoriesUseCaseImpl) List(filter *utils.QueryOptsBuilder) ([]Category, error) {
	categories, err := uc.repo.List(uc.db, filter)
	if err != nil {
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to list categories")
	}
	return categories, nil
}

func (uc *CategoriesUseCaseImpl) DeleteByID(id string) error {
	exists, err := uc.repo.Count(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete category")
	}

	if exists == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	err = uc.repo.DeleteByID(uc.db, id)

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete category")
	}

	return nil
}

func (uc *CategoriesUseCaseImpl) Count(filter *utils.QueryOptsBuilder) (int, error) {
	count, err := uc.repo.Count(uc.db, filter)

	if err != nil {
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count categories")
	}

	return count, nil
}

func (uc *CategoriesUseCaseImpl) ListCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) ([]CategoryAmountPerPeriod, error) {
	amounts, err := uc.repo.ListCategoryAmountPerPeriod(uc.db, period, filter)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to list category amounts per period")
	}
	return amounts, nil
}

func (uc *CategoriesUseCaseImpl) CountCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) (int, error) {
	count, err := uc.repo.CountCategoryAmountPerPeriod(uc.db, period, filter)

	if err != nil {
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count category amounts per period")
	}

	return count, nil
}

func (uc *CategoriesUseCaseImpl) Update(id string, payload UpdateCategoryDTO) (Category, error) {
	exists, err := uc.repo.Count(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
	}

	if exists == 0 {
		return Category{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	category, err := uc.repo.Update(uc.db, id, payload)

	if err != nil {
		return Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update category")
	}

	return category, nil
}
