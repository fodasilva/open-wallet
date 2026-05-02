// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *TransactionsRepoImpl) Update(ctx context.Context, db util.Executer, data UpdateTransactionDTO) error {
	filter := querybuilder.Get(ctx)
	var sets []string
	var values []interface{}
	if data.Name.Set {
		values = append(values, data.Name.Value)
		sets = append(sets, fmt.Sprintf("name = $%d", len(values)))
	}
	if data.Note.Set {
		values = append(values, data.Note.Value)
		sets = append(sets, fmt.Sprintf("description = $%d", len(values)))
	}
	if data.CategoryID.Set {
		values = append(values, data.CategoryID.Value)
		sets = append(sets, fmt.Sprintf("category_id = $%d", len(values)))
	}
	if data.RecurrenceID.Set {
		values = append(values, data.RecurrenceID.Value)
		sets = append(sets, fmt.Sprintf("recurrence_id = $%d", len(values)))
	}

	f := filter.ToSQL(len(values) + 1)
	sql := "UPDATE transactions SET " + strings.Join(sets, ", ") + " WHERE " + f.Where
	values = append(values, f.Args...)

	_, err := db.ExecContext(ctx, sql, values...)

	return err
}
