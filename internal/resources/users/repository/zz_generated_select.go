// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *UsersRepoImpl) Select(ctx context.Context, db utils.Executer) ([]User, error) {
	filter := querybuilder.FromContext(ctx)
	query := squirrel.Select("id", "name", "email", "avatar_url", "created_at", "username").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []User = []User{}
	for rows.Next() {
		var item User
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Email,
			&item.AvatarURL,
			&item.CreatedAt,
			&item.Username,
		)
		results = append(results, item)
	}

	return results, nil
}
