package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/felipe1496/open-wallet/internal/resources/auth/usecases"
	mockUsers "github.com/felipe1496/open-wallet/internal/resources/users/mocks"
	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	usersUseCases "github.com/felipe1496/open-wallet/internal/resources/users/usecases"
	"github.com/felipe1496/open-wallet/internal/services"
	mockServices "github.com/felipe1496/open-wallet/internal/services/mocks"
)

func TestAuthUseCase_LoginWithGoogle(t *testing.T) {
	t.Run("should not login if code is invalid", func(t *testing.T) {
		mockGoogle := new(mockServices.MockGoogleService)
		mockRepo := new(mockUsers.MockUsersRepo)
		usersUC := usersUseCases.NewUsersUseCases(mockRepo, nil)
		uc := usecases.NewAuthUseCases(mockGoogle, usersUC)

		mockGoogle.On("GetUserAccessToken", "invalid-code").Return(nil, services.FailedGoogleAuthenticationErr)

		_, err := uc.LoginWithGoogle("invalid-code")

		assert.ErrorIs(t, err, services.FailedGoogleAuthenticationErr)
		mockGoogle.AssertExpectations(t)
	})

	t.Run("should not login if google email is not verified", func(t *testing.T) {
		mockGoogle := new(mockServices.MockGoogleService)
		mockRepo := new(mockUsers.MockUsersRepo)
		usersUC := usersUseCases.NewUsersUseCases(mockRepo, nil)
		uc := usecases.NewAuthUseCases(mockGoogle, usersUC)

		accessToken := "valid-access-token"

		mockGoogle.
			On("GetUserAccessToken", "valid-code").
			Return(&accessToken, nil)

		emailVerified := false
		email := "test@gmail.com"
		mockGoogle.On("GetUserInfo", accessToken).Return(&services.GoogleUserInfo{
			Email:         &email,
			EmailVerified: &emailVerified,
		}, nil)

		_, err := uc.LoginWithGoogle("valid-code")

		assert.ErrorIs(t, err, usecases.GoogleEmailNotVerifiedErr)
		mockGoogle.AssertExpectations(t)
	})

	t.Run("should not login if google did not provide email", func(t *testing.T) {
		mockGoogle := new(mockServices.MockGoogleService)
		mockRepo := new(mockUsers.MockUsersRepo)
		usersUC := usersUseCases.NewUsersUseCases(mockRepo, nil)
		uc := usecases.NewAuthUseCases(mockGoogle, usersUC)

		accessToken := "valid-access-token"

		mockGoogle.
			On("GetUserAccessToken", "valid-code").
			Return(&accessToken, nil)

		emailVerified := true
		mockGoogle.On("GetUserInfo", accessToken).Return(&services.GoogleUserInfo{
			Email:         nil,
			EmailVerified: &emailVerified,
		}, nil)

		_, err := uc.LoginWithGoogle("valid-code")
		assert.ErrorIs(t, err, usecases.GoogleDintProvideEmailErr)
		mockGoogle.AssertExpectations(t)
	})

	t.Run("should return error if users usecase list fails", func(t *testing.T) {
		mockGoogleService := new(mockServices.MockGoogleService)
		mockRepo := new(mockUsers.MockUsersRepo)
		usersUseCase := usersUseCases.NewUsersUseCases(mockRepo, nil)
		uc := usecases.NewAuthUseCases(mockGoogleService, usersUseCase)

		emailVerified := true
		email := "valid@gmail.com"

		accessToken := "valid-access-token"

		mockGoogleService.
			On("GetUserAccessToken", "valid-code").
			Return(&accessToken, nil)

		mockGoogleService.
			On("GetUserInfo", accessToken).
			Return(&services.GoogleUserInfo{
				EmailVerified: &emailVerified,
				Email:         &email,
				Name:          "John Doe",
			}, nil)

		mockRepo.
			On("Select", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error"))

		_, err := uc.LoginWithGoogle("valid-code")

		assert.Contains(t, err.Error(), "failed to fetch users")

		mockGoogleService.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should_return_error_if_creating_user_fails", func(t *testing.T) {
		googleSvcMock := new(mockServices.MockGoogleService)
		usersRepoMock := new(mockUsers.MockUsersRepo)
		usersUseCase := usersUseCases.NewUsersUseCases(usersRepoMock, nil)
		authUseCase := usecases.NewAuthUseCases(googleSvcMock, usersUseCase)

		accessToken := "valid-access-token"

		googleSvcMock.
			On("GetUserAccessToken", "valid-code").
			Return(&accessToken, nil)

		verifiedEmail := true
		email := "newuser@gmail.com"
		googleSvcMock.
			On("GetUserInfo", accessToken).
			Return(&services.GoogleUserInfo{
				Email:         &email,
				EmailVerified: &verifiedEmail,
				Name:          "New User",
			}, nil)

		// List check in auth usecase
		usersRepoMock.
			On("Select", mock.Anything, mock.Anything).
			Return([]repository.User{}, nil).
			Once()

		// List checks (username and email) in users usecase Create
		usersRepoMock.
			On("Select", mock.Anything, mock.Anything).
			Return([]repository.User{}, nil).
			Times(2)

		usersRepoMock.
			On("Insert", mock.Anything, mock.Anything).
			Return(errors.New("db error")).
			Once()

		_, err := authUseCase.LoginWithGoogle("valid-code")

		assert.Contains(t, err.Error(), "failed to create user")
	})

	t.Run("should login successfully if user does not exist and is created", func(t *testing.T) {
		mockGoogle := new(mockServices.MockGoogleService)
		mockRepo := new(mockUsers.MockUsersRepo)
		usersUC := usersUseCases.NewUsersUseCases(mockRepo, nil)
		uc := usecases.NewAuthUseCases(mockGoogle, usersUC)

		accessToken := "valid-access-token"

		mockGoogle.
			On("GetUserAccessToken", "valid-code").
			Return(&accessToken, nil)

		emailVerified := true
		email := "newuser@gmail.com"
		mockGoogle.On("GetUserInfo", accessToken).Return(&services.GoogleUserInfo{
			Name:          "New User",
			Email:         &email,
			EmailVerified: &emailVerified,
		}, nil)

		// List check in auth usecase
		mockRepo.
			On("Select", mock.Anything, mock.Anything).
			Return([]repository.User{}, nil).Once()

		// List checks (username and email) in users usecase Create
		mockRepo.
			On("Select", mock.Anything, mock.Anything).
			Return([]repository.User{}, nil).Times(2)

		createdUser := repository.User{ID: "2", Name: "New User", Email: email}
		mockRepo.On("Insert", mock.Anything, mock.Anything).Return(nil).Once()

		// Always fetch after creation
		mockRepo.On("Select", mock.Anything, mock.Anything).Return([]repository.User{createdUser}, nil).Once()

		user, err := uc.LoginWithGoogle("valid-code")
		assert.NoError(t, err)
		assert.Equal(t, createdUser, user)
		mockGoogle.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
