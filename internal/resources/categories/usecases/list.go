package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/categories/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *CategoriesUseCasesImpl) List(filter *querybuilder.Builder) ([]repository.Category, error) {
	categories, err := uc.repo.Select(uc.db, filter)
	if err != nil {
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to list categories")
	}
	return categories, nil
}
