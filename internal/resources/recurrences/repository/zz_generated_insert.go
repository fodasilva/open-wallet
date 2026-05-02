// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/util"
)

func (r *RecurrencesRepoImpl) Insert(ctx context.Context, db util.Executer, data CreateRecurrenceDTO) error {
	var columns []string
	var values []interface{}
	var placeholders []string
	columns = append(columns, "id")
	values = append(values, data.ID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "user_id")
	values = append(values, data.UserID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "name")
	values = append(values, data.Name)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	if data.Note.Set {
		columns = append(columns, "note")
		values = append(values, data.Note.Value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	columns = append(columns, "amount")
	values = append(values, data.Amount)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "day_of_month")
	values = append(values, data.DayOfMonth)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	if data.CategoryID.Set {
		columns = append(columns, "category_id")
		values = append(values, data.CategoryID.Value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	columns = append(columns, "start_period")
	values = append(values, data.StartPeriod)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	if data.EndPeriod.Set {
		columns = append(columns, "end_period")
		values = append(values, data.EndPeriod.Value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}

	sql := fmt.Sprintf("INSERT INTO recurrences (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err := db.ExecContext(ctx, sql, values...)

	return err
}
