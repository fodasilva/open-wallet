// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *TransactionsRepoImpl) Select(db utils.Executer, filter *utils.QueryOptsBuilder) ([]Transaction, error) {
	query := squirrel.Select("id", "user_id", "category", "name", "description", "created_at", "category_id", "recurrence_id").
		From("transactions").
		PlaceholderFormat(squirrel.Dollar)

	query = utils.QueryOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []Transaction = []Transaction{}
	for rows.Next() {
		var item Transaction
		err = rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Type,
			&item.Name,
			&item.Description,
			&item.CreatedAt,
			&item.CategoryID,
			&item.RecurrenceID,
		)
		results = append(results, item)
	}

	return results, nil
}
func (r *EntriesRepoImpl) Select(db utils.Executer, filter *utils.QueryOptsBuilder) ([]ViewEntry, error) {
	query := squirrel.Select("id", "transaction_id", "name", "description", "amount", "period", "user_id", "category", "total_amount", "installment", "total_installments", "created_at", "reference_date", "category_id", "category_name", "category_color", "recurrence_id").
		From("v_entries").
		PlaceholderFormat(squirrel.Dollar)

	query = utils.QueryOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []ViewEntry = []ViewEntry{}
	for rows.Next() {
		var item ViewEntry
		err = rows.Scan(
			&item.ID,
			&item.TransactionID,
			&item.Name,
			&item.Description,
			&item.Amount,
			&item.Period,
			&item.UserID,
			&item.Type,
			&item.TotalAmount,
			&item.Installment,
			&item.TotalInstallments,
			&item.CreatedAt,
			&item.ReferenceDate,
			&item.CategoryID,
			&item.CategoryName,
			&item.CategoryColor,
			&item.RecurrenceID,
		)
		results = append(results, item)
	}

	return results, nil
}
