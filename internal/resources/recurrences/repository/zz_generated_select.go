// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

func (r *RecurrencesRepoImpl) Select(ctx context.Context, db utils.Executer) ([]Recurrence, error) {
	filter := querybuilder.FromContext(ctx)
	query := squirrel.Select("id", "user_id", "name", "note", "amount", "day_of_month", "category_id", "category_name", "category_color", "start_period", "end_period", "created_at").
		From("v_recurrences").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []Recurrence = []Recurrence{}
	for rows.Next() {
		var item Recurrence
		err = rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Note,
			&item.Amount,
			&item.DayOfMonth,
			&item.CategoryID,
			&item.CategoryName,
			&item.CategoryColor,
			&item.StartPeriod,
			&item.EndPeriod,
			&item.CreatedAt,
		)
		results = append(results, item)
	}

	return results, nil
}
