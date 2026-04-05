package repository

import (
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

type UsersRepo interface {
	Select(db utils.Executer, filter *querybuilder.Builder) ([]User, error)
	Insert(db utils.Executer, data CreateUserDTO) error
	Update(db utils.Executer, data UpdateUserDTO, filter *querybuilder.Builder) error
	Delete(db utils.Executer, filter *querybuilder.Builder) error
	Count(db utils.Executer, filter *querybuilder.Builder) (int, error)
}

type UsersRepoImpl struct{}

func NewUsersRepo() UsersRepo {
	return &UsersRepoImpl{}
}
