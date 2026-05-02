// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *TransactionsRepoImpl) Delete(ctx context.Context, db util.Executer) error {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "DELETE FROM transactions WHERE " + f.Where

	_, err := db.ExecContext(ctx, sql, f.Args...)
	return err
}

func (r *EntriesRepoImpl) Delete(ctx context.Context, db util.Executer) error {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "DELETE FROM entries WHERE " + f.Where

	_, err := db.ExecContext(ctx, sql, f.Args...)
	return err
}
