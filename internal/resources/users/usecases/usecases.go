package usecases

import (
	"context"
	"database/sql"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
)

type UsersUseCases interface {
	List(ctx context.Context) ([]repository.User, error)
	Create(ctx context.Context, input repository.CreateUserDTO) (repository.User, error)
}

type UsersUseCasesImpl struct {
	repo repository.UsersRepo
	db   *sql.DB
}

func NewUsersUseCases(repo repository.UsersRepo, db *sql.DB) UsersUseCases {
	return &UsersUseCasesImpl{repo: repo, db: db}
}
