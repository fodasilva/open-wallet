// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *UsersRepoImpl) Update(ctx context.Context, db utils.Executer, data UpdateUserDTO) error {
	filter := querybuilder.FromContext(ctx)
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

	query = querybuilder.ToUpdateSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
