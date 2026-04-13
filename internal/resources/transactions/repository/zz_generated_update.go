// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *TransactionsRepoImpl) Update(ctx context.Context, db utils.Executer, data UpdateTransactionDTO) error {
	filter := querybuilder.FromContext(ctx)
	query := squirrel.Update("transactions").
		PlaceholderFormat(squirrel.Dollar)

	if data.Name.Set {
		query = query.Set("name", data.Name.Value)
	}
	if data.Note.Set {
		query = query.Set("description", data.Note.Value)
	}
	if data.CategoryID.Set {
		query = query.Set("category_id", data.CategoryID.Value)
	}
	if data.RecurrenceID.Set {
		query = query.Set("recurrence_id", data.RecurrenceID.Value)
	}

	query = querybuilder.ToUpdateSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
