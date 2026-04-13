package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *CategoriesUseCasesImpl) CountCategoryAmountPerPeriod(ctx context.Context, period string) (int, error) {
	count, err := uc.repo.CountCategoryAmountPerPeriod(ctx, uc.db, period)

	if err != nil {
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count category amounts per period")
	}

	return count, nil
}
