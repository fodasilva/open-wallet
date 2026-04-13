// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *TransactionsRepoImpl) Insert(ctx context.Context, db utils.Executer, data CreateTransactionDTO) error {
	query := squirrel.Insert("transactions").
		PlaceholderFormat(squirrel.Dollar)

	var columns []string
	var values []interface{}
	columns = append(columns, "id")
	values = append(values, data.ID)
	columns = append(columns, "user_id")
	values = append(values, data.UserID)
	columns = append(columns, "category")
	values = append(values, data.Type)
	columns = append(columns, "name")
	values = append(values, data.Name)
	if data.Note.Set {
		columns = append(columns, "description")
		values = append(values, data.Note.Value)
	}
	if data.CategoryID.Set {
		columns = append(columns, "category_id")
		values = append(values, data.CategoryID.Value)
	}
	if data.RecurrenceID.Set {
		columns = append(columns, "recurrence_id")
		values = append(values, data.RecurrenceID.Value)
	}
	query = query.Columns(columns...).Values(values...)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
func (r *EntriesRepoImpl) Insert(ctx context.Context, db utils.Executer, data CreateEntryDTO) error {
	query := squirrel.Insert("entries").
		PlaceholderFormat(squirrel.Dollar)

	var columns []string
	var values []interface{}
	columns = append(columns, "id")
	values = append(values, data.ID)
	columns = append(columns, "transaction_id")
	values = append(values, data.TransactionID)
	columns = append(columns, "amount")
	values = append(values, data.Amount)
	columns = append(columns, "reference_date")
	values = append(values, data.ReferenceDate)
	query = query.Columns(columns...).Values(values...)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
