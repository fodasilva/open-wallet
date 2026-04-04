package usecases

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *UsersUseCasesImpl) List(filter *utils.QueryOptsBuilder) ([]repository.User, error) {
	users, err := uc.repo.Select(uc.db, filter)

	if err != nil {
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch users")
	}

	return users, nil
}
