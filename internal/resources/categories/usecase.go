package categories

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/oklog/ulid/v2"
)

type CategoriesUseCase interface {
	Create(payload repository.CreateCategoryDTO) (repository.Category, error)
	List(filter *utils.QueryOptsBuilder) ([]repository.Category, error)
	DeleteByID(id string) error
	Count(filter *utils.QueryOptsBuilder) (int, error)
	ListCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) ([]repository.CategoryAmountPerPeriod, error)
	CountCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) (int, error)
	Update(id string, payload repository.UpdateCategoryDTO) (repository.Category, error)
}

type CategoriesUseCaseImpl struct {
	repo repository.CategoriesRepo
	db   *sql.DB
}

func NewCategoriesUseCase(repo repository.CategoriesRepo, db *sql.DB) CategoriesUseCase {
	return &CategoriesUseCaseImpl{
		repo: repo,
		db:   db,
	}
}

func (uc *CategoriesUseCaseImpl) Create(payload repository.CreateCategoryDTO) (repository.Category, error) {
	payload.ID = ulid.Make().String()

	category, err := uc.repo.Insert(uc.db, payload)

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create category")
	}

	return category, nil
}

func (uc *CategoriesUseCaseImpl) List(filter *utils.QueryOptsBuilder) ([]repository.Category, error) {
	categories, err := uc.repo.Select(uc.db, filter)
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

	err = uc.repo.Delete(uc.db, utils.QueryOpts().And("id", "eq", id))

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

func (uc *CategoriesUseCaseImpl) ListCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) ([]repository.CategoryAmountPerPeriod, error) {
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

func (uc *CategoriesUseCaseImpl) Update(id string, payload repository.UpdateCategoryDTO) (repository.Category, error) {
	exists, err := uc.repo.Count(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
	}

	if exists == 0 {
		return repository.Category{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	err = uc.repo.Update(uc.db, payload, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update category")
	}

	category, err := uc.repo.Select(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to get updated category")
	}

	if len(category) == 0 {
		return repository.Category{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	return category[0], nil
}
