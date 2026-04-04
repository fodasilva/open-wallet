// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *UsersRepoImpl) Select(db utils.Executer, filter *utils.QueryOptsBuilder) ([]User, error) {
	query := squirrel.Select("id", "name", "email", "avatar_url", "created_at", "username").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	query = utils.QueryOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)

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
