package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *CategoriesUseCasesImpl) CountCategoryAmountPerPeriod(period string, filter *utils.QueryOptsBuilder) (int, error) {
	count, err := uc.repo.CountCategoryAmountPerPeriod(uc.db, period, filter)

	if err != nil {
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count category amounts per period")
	}

	return count, nil
}
