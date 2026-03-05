package transactions

import (
	"context"
	"errors"

	"github.com/felipe1496/open-wallet/internal/utils"
	"go.opentelemetry.io/otel"

	"github.com/Masterminds/squirrel"
	"github.com/oklog/ulid/v2"
)

type TransactionsRepo interface {
	CreateEntry(db utils.Executer, payload PersistEntryDTO) (Entry, error)
	CreateTransaction(db utils.Executer, payload CreateTransactionDTO) (Transaction, error)
	ListViewEntries(ctx context.Context, db utils.Executer, filter *utils.QueryOptsBuilder) ([]ViewEntry, error)
	CountViewEntries(ctx context.Context, db utils.Executer, filter *utils.QueryOptsBuilder) (int, error)
	DeleteTransactionById(db utils.Executer, id string) error
	ListTransactions(db utils.Executer, filter *utils.QueryOptsBuilder) ([]Transaction, error)
	UpdateTransaction(db utils.Executer, id string, payload UpdateTransactionDTO) (Transaction, error)
	DeleteEntry(db utils.Executer, filter *utils.QueryOptsBuilder) error
}

type TransactionsRepoImpl struct {
}

func NewTransactionsRepo(db utils.Executer) TransactionsRepo {
	return &TransactionsRepoImpl{}
}

func (r *TransactionsRepoImpl) CreateEntry(db utils.Executer, payload PersistEntryDTO) (Entry, error) {
	query, args, err := squirrel.Insert("entries").
		Columns("id", "transaction_id", "amount", "reference_date").
		Values(ulid.Make().String(), payload.TransactionID, payload.Amount, payload.ReferenceDate).
		Suffix("RETURNING id, transaction_id, amount, reference_date::text, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return Entry{}, err
	}

	var entry Entry
	err = db.QueryRow(query, args...).Scan(
		&entry.ID,
		&entry.TransactionID,
		&entry.Amount,
		&entry.ReferenceDate,
		&entry.CreatedAt,
	)

	return entry, err
}

func (r *TransactionsRepoImpl) CreateTransaction(db utils.Executer, payload CreateTransactionDTO) (Transaction, error) {
	query, args, err := squirrel.Insert("transactions").
		Columns("id", "user_id", "category", "name", "description", "category_id").
		Values(ulid.Make().String(), payload.UserID, payload.Type, payload.Name, &payload.Note, &payload.CategoryID).
		Suffix("RETURNING id, user_id, category, name, description, created_at, category_id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return Transaction{}, err
	}

	var transaction Transaction
	err = db.QueryRow(query, args...).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Type,
		&transaction.Name,
		&transaction.Description,
		&transaction.CreatedAt,
		&transaction.CategoryID,
	)
	return transaction, err
}

func (r *TransactionsRepoImpl) ListViewEntries(ctx context.Context, db utils.Executer, filter *utils.QueryOptsBuilder) ([]ViewEntry, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "TransactionsRepository.ListViewEntries")
	defer span.End()
	query := squirrel.Select("id", "transaction_id", "name", "description", "amount", "period", "user_id", "category", "total_amount", "installment", "total_installments", "created_at", "reference_date::text", "category_id", "category_name", "category_color").
		From("v_entries").
		PlaceholderFormat(squirrel.Dollar)

	query = utils.QueryOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	rows, err := db.QueryContext(ctx, sql, args...)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	defer rows.Close()

	var entries []ViewEntry = []ViewEntry{}
	for rows.Next() {
		var entry ViewEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.TransactionID,
			&entry.Name,
			&entry.Description,
			&entry.Amount,
			&entry.Period,
			&entry.UserID,
			&entry.Type,
			&entry.TotalAmount,
			&entry.Installment,
			&entry.TotalInstallments,
			&entry.CreatedAt,
			&entry.ReferenceDate,
			&entry.CategoryID,
			&entry.CategoryName,
			&entry.CategoryColor,
		); err != nil {
			span.RecordError(err)
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *TransactionsRepoImpl) CountViewEntries(ctx context.Context, db utils.Executer, filter *utils.QueryOptsBuilder) (int, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "TransactionsRepository.CountViewEntries")
	defer span.End()

	countQuery := squirrel.
		Select("COUNT(*)").
		From("v_entries").
		PlaceholderFormat(squirrel.Dollar)

	countQuery = utils.QueryOptsToSquirrel(
		countQuery,
		filter,
	)

	sql, args, err := countQuery.ToSql()
	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	var count int
	err = db.QueryRowContext(ctx, sql, args...).Scan(&count)

	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	return count, nil
}

func (r *TransactionsRepoImpl) DeleteTransactionById(db utils.Executer, id string) error {
	sql, args, err := squirrel.Delete("transactions").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}
func (r *TransactionsRepoImpl) ListTransactions(db utils.Executer, filter *utils.QueryOptsBuilder) ([]Transaction, error) {
	query := squirrel.Select("id", "user_id", "category", "name", "description", "created_at").
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

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Name,
			&transaction.Description,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionsRepoImpl) UpdateTransaction(db utils.Executer, id string, payload UpdateTransactionDTO) (Transaction, error) {
	if !utils.HasAtLeastOneField(payload) {
		return Transaction{}, errors.New("no fields to update")
	}

	query := squirrel.Update("transactions").Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, user_id, category, name, description, created_at, category_id").
		PlaceholderFormat(squirrel.Dollar)

	for _, field := range payload.Update {
		switch field {
		case "name":
			query = query.Set("name", payload.Name)
		case "note":
			query = query.Set("description", payload.Note)
		case "category_id":
			query = query.Set("category_id", payload.CategoryID)
		}
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return Transaction{}, err
	}

	var transaction Transaction
	err = db.QueryRow(sql, args...).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Type,
		&transaction.Name,
		&transaction.Description,
		&transaction.CreatedAt,
		&transaction.CategoryID,
	)

	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func (r *TransactionsRepoImpl) DeleteEntry(db utils.Executer, filter *utils.QueryOptsBuilder) error {
	query := squirrel.Delete("entries").PlaceholderFormat(squirrel.Dollar)

	if filter != nil {
		query = utils.DeleteOptsToSquirrel(query, filter)
	}

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}
