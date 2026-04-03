// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *CategoriesRepoImpl) Insert(db utils.Executer, data CreateCategoryDTO) error {
	query := squirrel.Insert("categories").
		PlaceholderFormat(squirrel.Dollar)

	var columns []string
	var values []interface{}
	columns = append(columns, "id")
	values = append(values, data.ID)
	columns = append(columns, "user_id")
	values = append(values, data.UserID)
	columns = append(columns, "name")
	values = append(values, data.Name)
	columns = append(columns, "color")
	values = append(values, data.Color)
	query = query.Columns(columns...).Values(values...)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}
