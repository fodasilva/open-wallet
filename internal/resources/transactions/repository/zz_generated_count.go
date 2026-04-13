// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *EntriesRepoImpl) Count(ctx context.Context, db utils.Executer) (int, error) {
	filter := querybuilder.FromContext(ctx)
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
