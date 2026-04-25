package handlers

import (
	"github.com/felipe1496/open-wallet/internal/resources/users/repository"
)

func MapUserResource(u repository.User) UserResource {
	return UserResource{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
	}
}
