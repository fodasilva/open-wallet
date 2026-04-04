package usecases

import (
	"database/sql"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type UsersUseCases interface {
	List(filter *utils.QueryOptsBuilder) ([]repository.User, error)
	Create(input repository.CreateUserDTO) (repository.User, error)
}

type UsersUseCasesImpl struct {
	repo repository.UsersRepo
	db   *sql.DB
}

func NewUsersUseCases(repo repository.UsersRepo, db *sql.DB) UsersUseCases {
	return &UsersUseCasesImpl{repo: repo, db: db}
}
