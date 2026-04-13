package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *CategoriesUseCasesImpl) List(ctx context.Context) ([]repository.Category, error) {
	categories, err := uc.repo.Select(ctx, uc.db)
	if err != nil {
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to list categories")
	}
	return categories, nil
}
