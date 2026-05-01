package usecases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oklog/ulid/v2"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/util/httputil"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (uc *UsersUseCasesImpl) Create(ctx context.Context, input repository.CreateUserDTO) (repository.User, error) {
	userAlreadyExists, err := uc.List(querybuilder.WithBuilder(ctx, querybuilder.New().And("username", "eq", input.Username)))
	if err == nil && len(userAlreadyExists) > 0 {
		return repository.User{}, httputil.NewHTTPError(http.StatusConflict, "user with this username already exists")
	}

	userAlreadyExists, err = uc.List(querybuilder.WithBuilder(ctx, querybuilder.New().And("email", "eq", input.Email)))
	if err == nil && len(userAlreadyExists) > 0 {
		return repository.User{}, httputil.NewHTTPError(http.StatusConflict, "user with this email already exists")
	}

	if input.ID == "" {
		input.ID = ulid.Make().String()
	}

	err = uc.repo.Insert(ctx, uc.db, input)
	if err != nil {
		return repository.User{}, httputil.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	createdCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("id", "eq", input.ID))
	created, err := uc.repo.Select(createdCtx, uc.db)
	if err != nil || len(created) == 0 {
		return repository.User{}, httputil.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch created user: %v", err))
	}

	return created[0], nil
}
