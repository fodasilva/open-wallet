// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/Masterminds/squirrel"

	"github.com/felipe1496/open-wallet/internal/util"
)

func (r *UsersRepoImpl) Insert(ctx context.Context, db util.Executer, data CreateUserDTO) error {
	query := squirrel.Insert("users").
		PlaceholderFormat(squirrel.Dollar)

	var columns []string
	var values []interface{}
	columns = append(columns, "id")
	values = append(values, data.ID)
	columns = append(columns, "name")
	values = append(values, data.Name)
	columns = append(columns, "email")
	values = append(values, data.Email)
	columns = append(columns, "avatar_url")
	values = append(values, data.AvatarURL)
	columns = append(columns, "username")
	values = append(values, data.Username)
	query = query.Columns(columns...).Values(values...)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
