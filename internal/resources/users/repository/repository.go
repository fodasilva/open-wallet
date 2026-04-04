package repository

import (
	"github.com/felipe1496/open-wallet/internal/utils"
)

type UsersRepo interface {
	Select(db utils.Executer, filter *utils.QueryOptsBuilder) ([]User, error)
	Insert(db utils.Executer, data CreateUserDTO) error
	Update(db utils.Executer, data UpdateUserDTO, filter *utils.QueryOptsBuilder) error
	Delete(db utils.Executer, filter *utils.QueryOptsBuilder) error
	Count(db utils.Executer, filter *utils.QueryOptsBuilder) (int, error)
}

type UsersRepoImpl struct{}

func NewUsersRepo() UsersRepo {
	return &UsersRepoImpl{}
}
