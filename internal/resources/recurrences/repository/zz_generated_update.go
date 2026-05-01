// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/Masterminds/squirrel"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *RecurrencesRepoImpl) Update(ctx context.Context, db util.Executer, data UpdateRecurrenceDTO) error {
	filter := querybuilder.Get(ctx)
	query := squirrel.Update("recurrences").
		PlaceholderFormat(squirrel.Dollar)

	if data.Name.Set {
		query = query.Set("name", data.Name.Value)
	}
	if data.Note.Set {
		query = query.Set("note", data.Note.Value)
	}
	if data.Amount.Set {
		query = query.Set("amount", data.Amount.Value)
	}
	if data.DayOfMonth.Set {
		query = query.Set("day_of_month", data.DayOfMonth.Value)
	}
	if data.CategoryID.Set {
		query = query.Set("category_id", data.CategoryID.Value)
	}
	if data.StartPeriod.Set {
		query = query.Set("start_period", data.StartPeriod.Value)
	}
	if data.EndPeriod.Set {
		query = query.Set("end_period", data.EndPeriod.Value)
	}

	query = querybuilder.ToUpdateSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sql, args...)

	return err
}
