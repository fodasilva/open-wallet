// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *UsersRepoImpl) Update(db utils.Executer, data UpdateUserDTO, filter *utils.QueryOptsBuilder) error {
	query := squirrel.Update("users").
		PlaceholderFormat(squirrel.Dollar)

	if data.Name.Set {
		query = query.Set("name", data.Name.Value)
	}
	if data.Email.Set {
		query = query.Set("email", data.Email.Value)
	}
	if data.AvatarURL.Set {
		query = query.Set("avatar_url", data.AvatarURL.Value)
	}
	if data.Username.Set {
		query = query.Set("username", data.Username.Value)
	}

	query = utils.UpdateOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}
