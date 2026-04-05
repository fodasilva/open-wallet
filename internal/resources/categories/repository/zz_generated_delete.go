// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *CategoriesRepoImpl) Delete(db utils.Executer, filter *querybuilder.Builder) error {
	query := squirrel.Delete("categories").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToDeleteSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}
