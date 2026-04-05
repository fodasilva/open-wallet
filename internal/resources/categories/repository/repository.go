package repository

import (
	"github.com/Masterminds/squirrel"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

// Repository interface. Make sure to only include methods
// that you defined with @method tags in types.go
type CategoriesRepo interface {
	Select(db utils.Executer, filter *querybuilder.Builder) ([]Category, error)
	Insert(db utils.Executer, data CreateCategoryDTO) error
	Update(db utils.Executer, data UpdateCategoryDTO, filter *querybuilder.Builder) error
	Delete(db utils.Executer, filter *querybuilder.Builder) error
	Count(db utils.Executer, filter *querybuilder.Builder) (int, error)
	CountCategoryAmountPerPeriod(db utils.Executer, period string, filter *querybuilder.Builder) (int, error)
	ListCategoryAmountPerPeriod(db utils.Executer, period string, filter *querybuilder.Builder) ([]CategoryAmountPerPeriod, error)
}

// Implementation struct. Name must match @name tag in types.go
type CategoriesRepoImpl struct {
}

func NewCategoriesRepo() CategoriesRepo {
	return &CategoriesRepoImpl{}
}

func (r *CategoriesRepoImpl) CountCategoryAmountPerPeriod(db utils.Executer, period string, filter *querybuilder.Builder) (int, error) {
	countQuery := squirrel.
		Select("COUNT(*)").
		From("fn_category_amount_per_period(?)").
		PlaceholderFormat(squirrel.Dollar)

	countQuery = querybuilder.ToSquirrel(countQuery, filter)

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

func (r *CategoriesRepoImpl) ListCategoryAmountPerPeriod(db utils.Executer, period string, filter *querybuilder.Builder) ([]CategoryAmountPerPeriod, error) {
	query := squirrel.
		Select("id", "user_id", "name", "color", "period", "total_amount").
		From("fn_category_amount_per_period(?)").
		PlaceholderFormat(squirrel.Dollar)

	query = querybuilder.ToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	args = append([]interface{}{period}, args...)

	rows, err := db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

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
