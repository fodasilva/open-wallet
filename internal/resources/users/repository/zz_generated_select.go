// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *UsersRepoImpl) Select(ctx context.Context, db util.Executer) ([]User, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "SELECT id, name, email, avatar_url, created_at, username FROM users WHERE " + f.Where + f.OrderBy + f.Limit + f.Offset

	rows, err := db.QueryContext(ctx, sql, f.Args...)

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
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}
