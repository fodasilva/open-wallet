package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *CategoriesUseCasesImpl) Count(ctx context.Context) (int, error) {
	count, err := uc.repo.Count(ctx, uc.db)

	if err != nil {
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count categories")
	}

	return count, nil
}
