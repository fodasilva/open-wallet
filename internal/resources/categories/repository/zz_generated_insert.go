// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *CategoriesRepoImpl) Insert(db utils.Executer, data CreateCategoryDTO) (Category, error) {
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

	query = query.Suffix("RETURNING id, user_id, name, color, created_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return Category{}, err
	}

	var result Category
	err = db.QueryRow(sql, args...).Scan(
		&result.ID,
		&result.UserID,
		&result.Name,
		&result.Color,
		&result.CreatedAt,
	)

	if err != nil {
		return Category{}, err
	}

	return result, nil
}
