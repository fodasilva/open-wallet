package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *CategoriesUseCasesImpl) Update(id string, userID string, payload repository.UpdateCategoryDTO) (repository.Category, error) {
	filter := utils.QueryOpts().And("id", "eq", id).And("user_id", "eq", userID)
	exists, err := uc.repo.Count(uc.db, filter)

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
	}

	if exists == 0 {
		return repository.Category{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	err = uc.repo.Update(uc.db, payload, filter)

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update category")
	}

	category, err := uc.repo.Select(uc.db, filter)

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to get updated category")
	}

	if len(category) == 0 {
		return repository.Category{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	return category[0], nil
}
