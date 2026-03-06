package recurrences

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
)

type RecurrencesRepo interface {
	Create(db utils.Executer, payload CreateRecurrenceDTO) (Recurrence, error)
	GetByID(ctx context.Context, db utils.Executer, id string) (Recurrence, error)
	List(ctx context.Context, db utils.Executer, filter *utils.QueryOptsBuilder) ([]Recurrence, error)
	DeleteByID(db utils.Executer, id string) error
	Count(ctx context.Context, db utils.Executer, filter *utils.QueryOptsBuilder) (int, error)
	Update(db utils.Executer, id string, payload UpdateRecurrenceDTO) (Recurrence, error)
}

type RecurrencesRepoImpl struct{}

func NewRecurrencesRepo() RecurrencesRepo {
	return &RecurrencesRepoImpl{}
}
func (r *RecurrencesRepoImpl) Create(db utils.Executer, payload CreateRecurrenceDTO) (Recurrence, error) {
	id := ulid.Make().String()
	query, args, err := squirrel.Insert("recurrences").
		Columns("id", "user_id", "name", "note", "amount", "day_of_month", "start_period", "end_period", "category_id").
		Values(id, payload.UserID, payload.Name, payload.Note, payload.Amount, payload.DayOfMonth, payload.StartPeriod, payload.EndPeriod, payload.CategoryID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return Recurrence{}, err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return Recurrence{}, err
	}

	return r.GetByID(context.Background(), db, id)
}

func (r *RecurrencesRepoImpl) GetByID(ctx context.Context, db utils.Executer, id string) (Recurrence, error) {
	query, args, err := squirrel.Select("id", "user_id", "name", "note", "amount", "day_of_month", "start_period", "end_period", "category_id", "created_at", "category_name", "category_color").
		From("v_recurrences").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return Recurrence{}, err
	}

	var rec Recurrence
	err = db.QueryRowContext(ctx, query, args...).Scan(
		&rec.ID,
		&rec.UserID,
		&rec.Name,
		&rec.Note,
		&rec.Amount,
		&rec.DayOfMonth,
		&rec.StartPeriod,
		&rec.EndPeriod,
		&rec.CategoryID,
		&rec.CreatedAt,
		&rec.CategoryName,
		&rec.CategoryColor,
	)
	return rec, err
}

func (r *RecurrencesRepoImpl) List(ctx context.Context, db utils.Executer, filter *utils.QueryOptsBuilder) ([]Recurrence, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "RecurrencesRepository.List")
	defer span.End()

	query := squirrel.Select("id", "user_id", "name", "note", "amount", "day_of_month", "start_period", "end_period", "category_id", "created_at", "category_name", "category_color").
		From("v_recurrences").
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

	var result []Recurrence = []Recurrence{}
	for rows.Next() {
		var rec Recurrence
		err = rows.Scan(
			&rec.ID,
			&rec.UserID,
			&rec.Name,
			&rec.Note,
			&rec.Amount,
			&rec.DayOfMonth,
			&rec.StartPeriod,
			&rec.EndPeriod,
			&rec.CategoryID,
			&rec.CreatedAt,
			&rec.CategoryName,
			&rec.CategoryColor,
		)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		result = append(result, rec)
	}

	return result, nil
}

func (r *RecurrencesRepoImpl) Count(ctx context.Context, db utils.Executer, filter *utils.QueryOptsBuilder) (int, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "RecurrencesRepository.Count")
	defer span.End()

	countQuery := squirrel.
		Select("COUNT(*)").
		From("recurrences").
		PlaceholderFormat(squirrel.Dollar)

	countQuery = utils.QueryOptsToSquirrel(countQuery, filter)

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

func (r *RecurrencesRepoImpl) DeleteByID(db utils.Executer, id string) error {
	sql, args, err := squirrel.Delete("recurrences").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}

func (r *RecurrencesRepoImpl) Update(db utils.Executer, id string, payload UpdateRecurrenceDTO) (Recurrence, error) {
	if len(payload.Update) == 0 {
		return Recurrence{}, errors.New("no fields to update")
	}

	query := squirrel.Update("recurrences").Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	for _, field := range payload.Update {
		switch field {
		case "name":
			query = query.Set("name", payload.Name)
		case "note":
			query = query.Set("note", payload.Note)
		case "amount":
			query = query.Set("amount", payload.Amount)
		case "day_of_month":
			query = query.Set("day_of_month", payload.DayOfMonth)
		case "start_period":
			query = query.Set("start_period", payload.StartPeriod)
		case "end_period":
			query = query.Set("end_period", payload.EndPeriod)
		case "category_id":
			query = query.Set("category_id", payload.CategoryID)
		}
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return Recurrence{}, err
	}

	_, err = db.Exec(sql, args...)
	if err != nil {
		return Recurrence{}, err
	}

	return r.GetByID(context.Background(), db, id)
}
