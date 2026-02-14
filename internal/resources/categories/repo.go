package categories

import (
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/Masterminds/squirrel"
	"github.com/oklog/ulid/v2"
)

type CategoriesRepo interface {
	Create(db utils.Executer, payload CreateCategoryDTO) (Category, error)
	List(db utils.Executer, filter *utils.QueryOptsBuilder) ([]Category, error)
	DeleteByID(db utils.Executer, id string) error
	Count(db utils.Executer, filter *utils.QueryOptsBuilder) (int, error)
	ListCategoryAmountPerPeriod(db utils.Executer, period string, filter *utils.QueryOptsBuilder) ([]CategoryAmountPerPeriod, error)
	CountCategoryAmountPerPeriod(db utils.Executer, period string, filter *utils.QueryOptsBuilder) (int, error)
	Update(db utils.Executer, id string, payload UpdateCategoryDTO) (Category, error)
}

type CategoriesRepoImpl struct {
}

func NewCategoriesRepo(db utils.Executer) CategoriesRepo {
	return &CategoriesRepoImpl{}
}

func (r *CategoriesRepoImpl) Create(db utils.Executer, payload CreateCategoryDTO) (Category, error) {
	query, args, err := squirrel.Insert("categories").
		Columns("id", "user_id", "name", "color").
		Values(ulid.Make().String(), payload.UserID, payload.Name, payload.Color).
		Suffix("RETURNING id, user_id, name, color, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return Category{}, err
	}

	var category Category
	err = db.QueryRow(query, args...).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Color,
		&category.CreatedAt,
	)
	return category, err
}

func (r *CategoriesRepoImpl) List(db utils.Executer, filter *utils.QueryOptsBuilder) ([]Category, error) {
	query := squirrel.Select("id", "user_id", "name", "color", "created_at").
		From("categories").
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

	var categories []Category = []Category{}
	for rows.Next() {
		var category Category
		err = rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Color,
			&category.CreatedAt,
		)
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoriesRepoImpl) Count(db utils.Executer, filter *utils.QueryOptsBuilder) (int, error) {
	countQuery := squirrel.
		Select("COUNT(*)").
		From("categories").
		PlaceholderFormat(squirrel.Dollar)

	countQuery = utils.QueryOptsToSquirrel(countQuery, filter)

	sql, args, err := countQuery.ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = db.QueryRow(sql, args...).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CategoriesRepoImpl) DeleteByID(db utils.Executer, id string) error {
	sql, args, err := squirrel.Delete("categories").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}

func (r *CategoriesRepoImpl) ListCategoryAmountPerPeriod(db utils.Executer, period string, filter *utils.QueryOptsBuilder) ([]CategoryAmountPerPeriod, error) {
	query := squirrel.
		Select("id", "user_id", "name", "color", "period", "total_amount").
		From("fn_category_amount_per_period(?)").
		PlaceholderFormat(squirrel.Dollar)

	query = utils.QueryOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	args = append([]interface{}{period}, args...)

	rows, err := db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]CategoryAmountPerPeriod, 0)
	for rows.Next() {

		var category CategoryAmountPerPeriod

		err = rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Color,
			&category.Period,
			&category.TotalAmount)
		if err != nil {
			return nil, err
		}

		result = append(result, category)
	}

	return result, nil
}

func (r *CategoriesRepoImpl) CountCategoryAmountPerPeriod(db utils.Executer, period string, filter *utils.QueryOptsBuilder) (int, error) {
	countQuery := squirrel.
		Select("COUNT(*)").
		From("fn_category_amount_per_period(?)").
		PlaceholderFormat(squirrel.Dollar)

	countQuery = utils.QueryOptsToSquirrel(countQuery, filter)

	sql, args, err := countQuery.ToSql()

	if err != nil {
		return 0, err
	}

	args = append([]interface{}{period}, args...)

	var count int
	err = db.QueryRow(sql, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CategoriesRepoImpl) Update(db utils.Executer, id string, payload UpdateCategoryDTO) (Category, error) {
	query := squirrel.Update("categories").Suffix("RETURNING id, user_id, name, color, created_at")

	if payload.Name != nil {
		query = query.Set("name", payload.Name)
	}

	if payload.Color != nil {
		query = query.Set("color", payload.Color)
	}

	sql, args, err := query.
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return Category{}, err
	}

	var category Category
	err = db.QueryRow(sql, args...).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Color,
		&category.CreatedAt,
	)

	return category, err
}
