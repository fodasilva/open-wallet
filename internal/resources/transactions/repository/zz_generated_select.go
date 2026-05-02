// Code generated. DO NOT EDIT.

package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *TransactionsRepoImpl) Select(ctx context.Context, db util.Executer) ([]Transaction, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "SELECT id, user_id, category, name, description, created_at, category_id, recurrence_id FROM transactions WHERE " + f.Where + f.OrderBy + f.Limit + f.Offset

	rows, err := db.QueryContext(ctx, sql, f.Args...)

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
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *EntriesRepoImpl) Select(ctx context.Context, db util.Executer) ([]ViewEntry, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "SELECT id, transaction_id, name, description, amount, period, user_id, category, total_amount, installment, total_installments, created_at, reference_date, category_id, category_name, category_color, recurrence_id FROM v_entries WHERE " + f.Where + f.OrderBy + f.Limit + f.Offset

	rows, err := db.QueryContext(ctx, sql, f.Args...)

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
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (r *SummariesRepoImpl) Select(ctx context.Context, db util.Executer) ([]ViewSummary, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(1)

	sql := "SELECT user_id, period, total_expense, total_income, total_balance FROM v_summaries WHERE " + f.Where + f.OrderBy + f.Limit + f.Offset

	rows, err := db.QueryContext(ctx, sql, f.Args...)

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
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}
