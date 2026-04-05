package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

type MockUsersRepo struct {
	mock.Mock
}

func (m *MockUsersRepo) Select(db utils.Executer, filter *querybuilder.Builder) ([]repository.User, error) {
	args := m.Called(db, filter)
	var res []repository.User
	if args.Get(0) != nil {
		res = args.Get(0).([]repository.User)
	}
	return res, args.Error(1)
}

func (m *MockUsersRepo) Insert(db utils.Executer, data repository.CreateUserDTO) error {
	args := m.Called(db, data)
	return args.Error(0)
}

func (m *MockUsersRepo) Update(db utils.Executer, data repository.UpdateUserDTO, filter *querybuilder.Builder) error {
	args := m.Called(db, data, filter)
	return args.Error(0)
}

func (m *MockUsersRepo) Delete(db utils.Executer, filter *querybuilder.Builder) error {
	args := m.Called(db, filter)
	return args.Error(0)
}

func (m *MockUsersRepo) Count(db utils.Executer, filter *querybuilder.Builder) (int, error) {
	args := m.Called(db, filter)
	return args.Int(0), args.Error(1)
}
