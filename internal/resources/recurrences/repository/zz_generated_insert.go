// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *RecurrencesRepoImpl) Insert(db utils.Executer, data CreateRecurrenceDTO) error {
	query := squirrel.Insert("recurrences").
		PlaceholderFormat(squirrel.Dollar)

	var columns []string
	var values []interface{}
	columns = append(columns, "id")
	values = append(values, data.ID)
	columns = append(columns, "user_id")
	values = append(values, data.UserID)
	columns = append(columns, "name")
	values = append(values, data.Name)
	if data.Note.Set {
		columns = append(columns, "note")
		values = append(values, data.Note.Value)
	}
	columns = append(columns, "amount")
	values = append(values, data.Amount)
	columns = append(columns, "day_of_month")
	values = append(values, data.DayOfMonth)
	if data.CategoryID.Set {
		columns = append(columns, "category_id")
		values = append(values, data.CategoryID.Value)
	}
	columns = append(columns, "start_period")
	values = append(values, data.StartPeriod)
	if data.EndPeriod.Set {
		columns = append(columns, "end_period")
		values = append(values, data.EndPeriod.Value)
	}
	query = query.Columns(columns...).Values(values...)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}
