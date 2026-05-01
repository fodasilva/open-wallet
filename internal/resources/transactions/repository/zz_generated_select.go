// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *TransactionsRepoImpl) Select(ctx context.Context, db util.Executer) ([]Transaction, error) {
	filter := querybuilder.Get(ctx)
	query := squirrel.Select("id", "user_id", "category", "name", "description", "created_at", "category_id", "recurrence_id").
		From("transactions").
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

func (r *EntriesRepoImpl) Select(ctx context.Context, db util.Executer) ([]ViewEntry, error) {
	filter := querybuilder.Get(ctx)
	query := squirrel.Select("id", "transaction_id", "name", "description", "amount", "period", "user_id", "category", "total_amount", "installment", "total_installments", "created_at", "reference_date", "category_id", "category_name", "category_color", "recurrence_id").
		From("v_entries").
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

func (r *SummariesRepoImpl) Select(ctx context.Context, db util.Executer) ([]ViewSummary, error) {
	filter := querybuilder.Get(ctx)
	query := squirrel.Select("user_id", "period", "total_expense", "total_income", "total_balance").
		From("v_summaries").
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

	var results []ViewSummary = []ViewSummary{}
	for rows.Next() {
		var item ViewSummary
		err = rows.Scan(
			&item.UserID,
			&item.Period,
			&item.TotalExpense,
			&item.TotalIncome,
			&item.TotalBalance,
		)
		results = append(results, item)
	}

	return results, nil
}
