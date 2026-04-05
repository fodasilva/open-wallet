package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *CategoriesUseCasesImpl) DeleteByID(id string, userID string) error {
	filter := querybuilder.New().And("id", "eq", id).And("user_id", "eq", userID)
	exists, err := uc.repo.Count(uc.db, filter)

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete category")
	}

	if exists == 0 {
		return utils.NewHTTPError(http.StatusNotFound, "category not found")
	}

	err = uc.repo.Delete(uc.db, filter)

	if err != nil {
		return utils.NewHTTPError(http.StatusInternalServerError, "failed to delete category")
	}

	return nil
}
