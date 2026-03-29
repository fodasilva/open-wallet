// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *CategoriesRepoImpl) Update(db utils.Executer, data UpdateCategoryDTO, filter *utils.QueryOptsBuilder) error {
	query := squirrel.Update("categories").
		PlaceholderFormat(squirrel.Dollar)

	if data.Name.Set {
		query = query.Set("name", data.Name.Value)
	}
	if data.Color.Set {
		query = query.Set("color", data.Color.Value)
	}

	query = utils.UpdateOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}
