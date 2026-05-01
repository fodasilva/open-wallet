package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/util"
)

func (uc *CategoriesUseCasesImpl) Count(ctx context.Context) (int, error) {
	count, err := uc.repo.Count(ctx, uc.db)

	if err != nil {
		return 0, util.NewHTTPError(http.StatusInternalServerError, "failed to count categories")
	}

	return count, nil
}
