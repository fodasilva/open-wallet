package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *CategoriesUseCasesImpl) DeleteByID(ctx context.Context, id string, userID string) error {
	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", id).And("user_id", "eq", userID))
	exists, err := uc.repo.Count(filterCtx, uc.db)

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete category")
	}

	if exists == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	err = uc.repo.Delete(filterCtx, uc.db)

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete category")
	}

	return nil
}
