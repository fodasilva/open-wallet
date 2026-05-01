package repository

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/util"
)

// @gen_repo
// @table: users
// @entity: User
// @name: UsersRepoImpl
// @method: Select | fields: id:ID, name:Name, email:Email, avatar_url:AvatarURL, created_at:CreatedAt, username:Username
// @method: Insert | fields: id:ID, name:Name, email:Email, avatar_url:AvatarURL, username:Username | payload: CreateUserDTO
// @method: Update | fields: name:Name?, email:Email?, avatar_url:AvatarURL?, username:Username? | payload: UpdateUserDTO
// @method: Delete
// @method: Count

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
}

type CreateUserDTO struct {
	ID        string
	Name      string
	Email     string
	AvatarURL string
	Username  string
}

type UpdateUserDTO struct {
	Name      util.OptionalNullable[string]
	Email     util.OptionalNullable[string]
	AvatarURL util.OptionalNullable[string]
	Username  util.OptionalNullable[string]
}
