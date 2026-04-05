package usecases

import (
	"fmt"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *CategoriesUseCasesImpl) ListCategoryAmountPerPeriod(period string, filter *querybuilder.Builder) ([]repository.CategoryAmountPerPeriod, error) {
	amounts, err := uc.repo.ListCategoryAmountPerPeriod(uc.db, period, filter)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to list category amounts per period")
	}
	return amounts, nil
}
