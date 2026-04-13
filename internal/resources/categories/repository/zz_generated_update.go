// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *CategoriesRepoImpl) Update(ctx context.Context, db utils.Executer, data UpdateCategoryDTO) error {
	filter := querybuilder.FromContext(ctx)
	query := squirrel.Update("categories").
		PlaceholderFormat(squirrel.Dollar)

	if data.Name.Set {
		query = query.Set("name", data.Name.Value)
	}
	if data.Color.Set {
		query = query.Set("color", data.Color.Value)
	}

	query = querybuilder.ToUpdateSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
