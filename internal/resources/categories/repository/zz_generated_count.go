// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *CategoriesRepoImpl) Count(ctx context.Context, db util.Executer) (int, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "SELECT COUNT(*) FROM categories WHERE " + f.Where

	var count int
	err := db.QueryRowContext(ctx, sql, f.Args...).Scan(&count)
	return count, err
}
