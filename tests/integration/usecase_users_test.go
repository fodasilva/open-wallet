package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/felipe1496/open-wallet/internal/resources/users/mocks"
	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	usersUseCases "github.com/felipe1496/open-wallet/internal/resources/users/usecases"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func TestUsersUseCase_List(t *testing.T) {
	t.Run("should list users successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := usersUseCases.NewUsersUseCases(mockRepo, nil)

		expectedUsers := []repository.User{
			{ID: "1", Username: "alice", Name: "Alice", Email: "alice@gmail.com"},
			{ID: "2", Username: "alice2", Name: "Alice2", Email: "alice2@gmail.com"},
		}

		// Repo returns users successfully
		mockRepo.
			On("Select", mock.Anything, mock.Anything).
			Return(expectedUsers, nil)

		result, err := uc.List(querybuilder.New().And("username", "eq", "alice"))

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := usersUseCases.NewUsersUseCases(mockRepo, nil)

		mockRepo.
			On("Select", mock.Anything, mock.Anything).
			Return(nil, errors.New("db exploded"))

		result, err := uc.List(querybuilder.New().And("email", "eq", "john@gmail.com"))

		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to fetch users")
		mockRepo.AssertExpectations(t)
	})
}

func TestUsersUseCase_Create(t *testing.T) {
	t.Run("should return error if username is already taken", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := usersUseCases.NewUsersUseCases(mockRepo, nil)

		input := repository.CreateUserDTO{Username: "johndoethegreat", Name: "John", Email: "john@gmail.com"}

		mockRepo.
			On("Select", mock.Anything, mock.MatchedBy(func(filter *querybuilder.Builder) bool {
				// Very simple check to identify which Select call this is
				return filter != nil
			})).
			Return([]repository.User{
				{
					ID:       "1",
					Username: "johndoethegreat",
					Name:     "Urek",
					Email:    "urek@gmail.com",
				},
			}, nil).Once()

		result, err := uc.Create(input)

		assert.Equal(t, repository.User{}, result)
		assert.Contains(t, err.Error(), "user with this username already exists")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if email already exists", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := usersUseCases.NewUsersUseCases(mockRepo, nil)

		input := repository.CreateUserDTO{
			Username: "johndoethegreat",
			Name:     "John",
			Email:    "john@gmail.com",
		}

		// First Select: username check
		mockRepo.
			On("Select", mock.Anything, mock.Anything).
			Return([]repository.User{}, nil).Once()

		// Second Select: email check
		mockRepo.
			On("Select", mock.Anything, mock.Anything).
			Return([]repository.User{
				{
					ID:       "1",
					Username: "rolling_stone",
					Name:     "Urek",
					Email:    "john@gmail.com",
				},
			}, nil).Once()

		result, err := uc.Create(input)

		assert.Equal(t, repository.User{}, result)
		assert.Contains(t, err.Error(), "user with this email already exists")

		mockRepo.AssertExpectations(t)
	})

	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := usersUseCases.NewUsersUseCases(mockRepo, nil)

		input := repository.CreateUserDTO{
			Username: "johndoethegreat",
			Name:     "John",
			Email:    "john@gmail.com",
		}

		expectedUser := repository.User{
			ID:       "123",
			Username: input.Username,
			Name:     input.Name,
			Email:    input.Email,
		}

		// List checks (username and email)
		mockRepo.On("Select", mock.Anything, mock.Anything).Return([]repository.User{}, nil).Times(2)

		mockRepo.On("Insert", mock.Anything, mock.Anything).Return(nil).Once()

		// Always fetch
		mockRepo.On("Select", mock.Anything, mock.Anything).Return([]repository.User{expectedUser}, nil).Once()

		result, err := uc.Create(input)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})
}
