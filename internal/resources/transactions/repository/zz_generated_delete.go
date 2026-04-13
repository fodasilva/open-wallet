// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *TransactionsRepoImpl) Delete(ctx context.Context, db utils.Executer) error {
	filter := querybuilder.FromContext(ctx)
	query := squirrel.Delete("transactions").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToDeleteSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
func (r *EntriesRepoImpl) Delete(ctx context.Context, db utils.Executer) error {
	filter := querybuilder.FromContext(ctx)
	query := squirrel.Delete("entries").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToDeleteSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
