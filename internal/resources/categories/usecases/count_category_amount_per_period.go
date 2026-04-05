package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *CategoriesUseCasesImpl) CountCategoryAmountPerPeriod(period string, filter *querybuilder.Builder) (int, error) {
	count, err := uc.repo.CountCategoryAmountPerPeriod(uc.db, period, filter)

	if err != nil {
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count category amounts per period")
	}

	return count, nil
}
