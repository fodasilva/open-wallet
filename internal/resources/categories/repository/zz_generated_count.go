// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *CategoriesRepoImpl) Count(db utils.Executer, filter *utils.QueryOptsBuilder) (int, error) {
	query := squirrel.Select("COUNT(*)").
		From("categories").
		PlaceholderFormat(squirrel.Dollar)

	query = utils.QueryOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return 0, err
	}

	var count int
	err = db.QueryRow(sql, args...).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
