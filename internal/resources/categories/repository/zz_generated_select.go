// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *CategoriesRepoImpl) Select(ctx context.Context, db util.Executer) ([]Category, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "SELECT id, user_id, name, color, created_at FROM categories WHERE " + f.Where + f.OrderBy + f.Limit + f.Offset

	rows, err := db.QueryContext(ctx, sql, f.Args...)

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
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}
