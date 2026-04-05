package usecases

import (
	"fmt"
	"net/http"

	"github.com/oklog/ulid/v2"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *UsersUseCasesImpl) Create(input repository.CreateUserDTO) (repository.User, error) {
	userAlreadyExists, err := uc.List(querybuilder.New().And("username", "eq", input.Username))
	if err == nil && len(userAlreadyExists) > 0 {
		return repository.User{}, utils.NewHTTPError(http.StatusConflict, "user with this username already exists")
	}

	userAlreadyExists, err = uc.List(querybuilder.New().And("email", "eq", input.Email))
	if err == nil && len(userAlreadyExists) > 0 {
		return repository.User{}, utils.NewHTTPError(http.StatusConflict, "user with this email already exists")
	}

	if input.ID == "" {
		input.ID = ulid.Make().String()
	}

	err = uc.repo.Insert(uc.db, input)
	if err != nil {
		return repository.User{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	created, err := uc.repo.Select(uc.db, querybuilder.New().And("id", "eq", input.ID))
	if err != nil || len(created) == 0 {
		return repository.User{}, utils.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch created user: %v", err))
	}

	return created[0], nil
}
