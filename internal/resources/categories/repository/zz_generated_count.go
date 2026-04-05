// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *CategoriesRepoImpl) Count(db utils.Executer, filter *querybuilder.Builder) (int, error) {
	query := squirrel.Select("COUNT(*)").
		From("categories").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToSquirrel(query, filter)

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
