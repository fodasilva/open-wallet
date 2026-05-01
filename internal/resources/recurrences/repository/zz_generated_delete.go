// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/Masterminds/squirrel"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *RecurrencesRepoImpl) Delete(ctx context.Context, db util.Executer) error {
	filter := querybuilder.Get(ctx)
	query := squirrel.Delete("recurrences").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToDeleteSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
