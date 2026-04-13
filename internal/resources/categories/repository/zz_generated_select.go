// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *CategoriesRepoImpl) Select(ctx context.Context, db utils.Executer) ([]Category, error) {
	filter := querybuilder.FromContext(ctx)
	query := squirrel.Select("id", "user_id", "name", "color", "created_at").
		From("categories").
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

	var results []Category = []Category{}
	for rows.Next() {
		var item Category
		err = rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Color,
			&item.CreatedAt,
		)
		results = append(results, item)
	}

	return results, nil
}
