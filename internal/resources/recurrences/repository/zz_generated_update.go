// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *RecurrencesRepoImpl) Update(ctx context.Context, db util.Executer, data UpdateRecurrenceDTO) error {
	filter := querybuilder.Get(ctx)
	var sets []string
	var values []interface{}
	if data.Name.Set {
		values = append(values, data.Name.Value)
		sets = append(sets, fmt.Sprintf("name = $%d", len(values)))
	}
	if data.Note.Set {
		values = append(values, data.Note.Value)
		sets = append(sets, fmt.Sprintf("note = $%d", len(values)))
	}
	if data.Amount.Set {
		values = append(values, data.Amount.Value)
		sets = append(sets, fmt.Sprintf("amount = $%d", len(values)))
	}
	if data.DayOfMonth.Set {
		values = append(values, data.DayOfMonth.Value)
		sets = append(sets, fmt.Sprintf("day_of_month = $%d", len(values)))
	}
	if data.CategoryID.Set {
		values = append(values, data.CategoryID.Value)
		sets = append(sets, fmt.Sprintf("category_id = $%d", len(values)))
	}
	if data.StartPeriod.Set {
		values = append(values, data.StartPeriod.Value)
		sets = append(sets, fmt.Sprintf("start_period = $%d", len(values)))
	}
	if data.EndPeriod.Set {
		values = append(values, data.EndPeriod.Value)
		sets = append(sets, fmt.Sprintf("end_period = $%d", len(values)))
	}

	f := filter.ToSQL(len(values) + 1)
	sql := "UPDATE recurrences SET " + strings.Join(sets, ", ") + " WHERE " + f.Where
	values = append(values, f.Args...)

	_, err := db.ExecContext(ctx, sql, values...)

	return err
}
