package usecases

import (
	"context"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (uc *UsersUseCasesImpl) List(ctx context.Context) ([]repository.User, error) {
	users, err := uc.repo.Select(ctx, uc.db)

	if err != nil {
		return nil, utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch users")
	}

	return users, nil
}
