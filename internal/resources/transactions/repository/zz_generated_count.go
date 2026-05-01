// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/Masterminds/squirrel"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *EntriesRepoImpl) Count(ctx context.Context, db util.Executer) (int, error) {
	filter := querybuilder.Get(ctx)
	query := squirrel.Select("COUNT(*)").
		From("v_entries").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return 0, err
	}

	var count int
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
