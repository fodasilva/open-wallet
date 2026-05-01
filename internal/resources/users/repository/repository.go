package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
)

type UsersRepo interface {
	Select(ctx context.Context, db util.Executer) ([]User, error)
	Insert(ctx context.Context, db util.Executer, data CreateUserDTO) error
	Update(ctx context.Context, db util.Executer, data UpdateUserDTO) error
	Delete(ctx context.Context, db util.Executer) error
	Count(ctx context.Context, db util.Executer) (int, error)
}

type UsersRepoImpl struct{}

func NewUsersRepo() UsersRepo {
	return &UsersRepoImpl{}
}
