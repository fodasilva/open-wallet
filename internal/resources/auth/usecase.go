package auth

import (
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/resources/users"
	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/services"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/oklog/ulid/v2"
)

type AuthUseCase interface {
	LoginWithGoogle(code string) (repository.User, error)
}

type AuthUseCaseImpl struct {
	googleService services.GoogleService
	usersUseCase  users.UsersUseCase
}

func NewAuthUseCase(googleService services.GoogleService, usersUseCase users.UsersUseCase) AuthUseCase {
	return &AuthUseCaseImpl{
		googleService: googleService,
		usersUseCase:  usersUseCase,
	}
}

func (uc *AuthUseCaseImpl) LoginWithGoogle(code string) (repository.User, error) {
	userAccessToken, err := uc.googleService.GetUserAccessToken(code)

	if err != nil {
		return repository.User{}, err
	}

	userInfo, err := uc.googleService.GetUserInfo(*userAccessToken)

	if err != nil {
		return repository.User{}, err
	}

	if !*userInfo.EmailVerified {
		return repository.User{}, GoogleEmailNotVerifiedErr
	}

	if userInfo.Email == nil {
		return repository.User{}, GoogleDintProvideEmailErr
	}

	userExists, err := uc.usersUseCase.List(utils.QueryOpts().And("email", "eq", *userInfo.Email))

	if err != nil {
		return repository.User{}, err
	}

	var userRes repository.User

	if len(userExists) == 0 {
		createUserInput := repository.CreateUserDTO{
			Name: userInfo.Name,
		}

		createUserInput.Email = *userInfo.Email

		if userInfo.Picture != nil {
			createUserInput.AvatarURL = *userInfo.Picture
		}

		createUserInput.Username = fmt.Sprintf("%s_%s", strings.ToLower(strings.ReplaceAll(userInfo.Name, " ", "_")), ulid.Make().String())

		createdUser, err := uc.usersUseCase.Create(createUserInput)

		if err != nil {
			return repository.User{}, err
		}

		userRes = createdUser
	} else {
		userRes = userExists[0]
	}

	return userRes, nil
}
