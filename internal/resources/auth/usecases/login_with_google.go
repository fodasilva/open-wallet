package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"

	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (uc *AuthUseCasesImpl) LoginWithGoogle(ctx context.Context, code string) (repository.User, error) {
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

	filterCtx := querybuilder.WithBuilder(ctx, querybuilder.New().And("email", "eq", *userInfo.Email))
	userExists, err := uc.usersUseCase.List(filterCtx)

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

		createdUser, err := uc.usersUseCase.Create(ctx, createUserInput)

		if err != nil {
			return repository.User{}, err
		}

		userRes = createdUser
	} else {
		userRes = userExists[0]
	}

	return userRes, nil
}
