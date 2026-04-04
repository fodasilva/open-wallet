package users

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/oklog/ulid/v2"
)

type UsersUseCase interface {
	List(filter *utils.QueryOptsBuilder) ([]repository.User, error)
	Create(input repository.CreateUserDTO) (repository.User, error)
}

type UsersUseCaseImpl struct {
	repo repository.UsersRepo
	db   *sql.DB
}

func NewUsersUseCase(repo repository.UsersRepo, db *sql.DB) UsersUseCase {
	return &UsersUseCaseImpl{repo: repo, db: db}
}

func (uc *UsersUseCaseImpl) List(filter *utils.QueryOptsBuilder) ([]repository.User, error) {
	users, err := uc.repo.Select(uc.db, filter)

	if err != nil {
		return nil, FailedToFetchUsersError
	}

	return users, nil
}

func (uc *UsersUseCaseImpl) Create(input repository.CreateUserDTO) (repository.User, error) {
	// Check if user already exists by username
	userAlreadyExists, err := uc.List(utils.QueryOpts().And("username", "eq", input.Username))
	if err == nil && len(userAlreadyExists) > 0 {
		return repository.User{}, UsernameAlreadyExists
	}

	// Check if user already exists by email
	userAlreadyExists, err = uc.List(utils.QueryOpts().And("email", "eq", input.Email))
	if err == nil && len(userAlreadyExists) > 0 {
		return repository.User{}, EmailAlreadyExists
	}

	if input.ID == "" {
		input.ID = ulid.Make().String()
	}

	err = uc.repo.Insert(uc.db, input)
	if err != nil {
		return repository.User{}, FailedToCreateUserError
	}

	// Always fetch after creation
	created, err := uc.repo.Select(uc.db, utils.QueryOpts().And("id", "eq", input.ID))
	if err != nil || len(created) == 0 {
		return repository.User{}, utils.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch created user: %v", err))
	}

	return created[0], nil
}
