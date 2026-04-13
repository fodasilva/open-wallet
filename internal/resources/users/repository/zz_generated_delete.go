// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *UsersRepoImpl) Delete(ctx context.Context, db utils.Executer) error {
	filter := querybuilder.FromContext(ctx)
	query := squirrel.Delete("users").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToDeleteSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
