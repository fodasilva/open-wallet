package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *CategoriesUseCasesImpl) Count(filter *querybuilder.Builder) (int, error) {
	count, err := uc.repo.Count(uc.db, filter)

	if err != nil {
		return 0, utils.NewHTTPError(http.StatusInternalServerError, "failed to count categories")
	}

	return count, nil
}
