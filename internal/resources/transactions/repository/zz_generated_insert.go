// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/util"
)

func (r *TransactionsRepoImpl) Insert(ctx context.Context, db util.Executer, data CreateTransactionDTO) error {
	var columns []string
	var values []interface{}
	var placeholders []string
	columns = append(columns, "id")
	values = append(values, data.ID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "user_id")
	values = append(values, data.UserID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "category")
	values = append(values, data.Type)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "name")
	values = append(values, data.Name)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	if data.Note.Set {
		columns = append(columns, "description")
		values = append(values, data.Note.Value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if data.CategoryID.Set {
		columns = append(columns, "category_id")
		values = append(values, data.CategoryID.Value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if data.RecurrenceID.Set {
		columns = append(columns, "recurrence_id")
		values = append(values, data.RecurrenceID.Value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}

	sql := fmt.Sprintf("INSERT INTO transactions (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err := db.ExecContext(ctx, sql, values...)

	return err
}

func (r *EntriesRepoImpl) Insert(ctx context.Context, db util.Executer, data CreateEntryDTO) error {
	var columns []string
	var values []interface{}
	var placeholders []string
	columns = append(columns, "id")
	values = append(values, data.ID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "transaction_id")
	values = append(values, data.TransactionID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "amount")
	values = append(values, data.Amount)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "reference_date")
	values = append(values, data.ReferenceDate)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))

	sql := fmt.Sprintf("INSERT INTO entries (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err := db.ExecContext(ctx, sql, values...)

	return err
}
