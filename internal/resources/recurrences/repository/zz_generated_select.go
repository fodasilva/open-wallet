// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *RecurrencesRepoImpl) Select(ctx context.Context, db util.Executer) ([]Recurrence, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "SELECT id, user_id, name, note, amount, day_of_month, category_id, category_name, category_color, start_period, end_period, created_at FROM v_recurrences WHERE " + f.Where + f.OrderBy + f.Limit + f.Offset

	rows, err := db.QueryContext(ctx, sql, f.Args...)

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
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}
