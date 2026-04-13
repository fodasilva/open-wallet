package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *CategoriesUseCasesImpl) Update(ctx context.Context, id string, userID string, payload repository.UpdateCategoryDTO) (repository.Category, error) {
	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", id).And("user_id", "eq", userID))
	exists, err := uc.repo.Count(filterCtx, uc.db)

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
	}

	if exists == 0 {
		return repository.Category{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	err = uc.repo.Update(filterCtx, uc.db, payload)

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update category")
	}

	category, err := uc.repo.Select(filterCtx, uc.db)

	if err != nil {
		return repository.Category{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to get updated category")
	}

	if len(category) == 0 {
		return repository.Category{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	return category[0], nil
}
