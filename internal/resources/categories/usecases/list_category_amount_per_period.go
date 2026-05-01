package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/util"
)

func (uc *CategoriesUseCasesImpl) ListCategoryAmountPerPeriod(ctx context.Context, period string) ([]repository.CategoryAmountPerPeriod, error) {
	amounts, err := uc.repo.ListCategoryAmountPerPeriod(ctx, uc.db, period)
	if err != nil {
		return nil, util.NewHTTPError(http.StatusInternalServerError, "failed to list category amounts per period")
	}
	return amounts, nil
}
