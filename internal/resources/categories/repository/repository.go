package repository

import (
	"context"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

// Repository interface. Make sure to only include methods
// that you defined with @method tags in types.go
type CategoriesRepo interface {
	Select(ctx context.Context, db util.Executer) ([]Category, error)
	Insert(ctx context.Context, db util.Executer, data CreateCategoryDTO) error
	Update(ctx context.Context, db util.Executer, data UpdateCategoryDTO) error
	Delete(ctx context.Context, db util.Executer) error
	Count(ctx context.Context, db util.Executer) (int, error)
	CountCategoryAmountPerPeriod(ctx context.Context, db util.Executer, period string) (int, error)
	ListCategoryAmountPerPeriod(ctx context.Context, db util.Executer, period string) ([]CategoryAmountPerPeriod, error)
}

// Implementation struct. Name must match @name tag in types.go
type CategoriesRepoImpl struct {
}

func NewCategoriesRepo() CategoriesRepo {
	return &CategoriesRepoImpl{}
}

func (r *CategoriesRepoImpl) CountCategoryAmountPerPeriod(ctx context.Context, db util.Executer, period string) (int, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(2) // Starts at $2 because period is $1

	sql := "SELECT COUNT(*) FROM fn_category_amount_per_period($1) WHERE " + f.Where
	args := append([]any{period}, f.Args...)

	var count int
	err := db.QueryRowContext(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CategoriesRepoImpl) ListCategoryAmountPerPeriod(ctx context.Context, db util.Executer, period string) ([]CategoryAmountPerPeriod, error) {
	filter := querybuilder.Get(ctx)
	f := filter.ToSQL(2) // Starts at $2 because period is $1

	sql := "SELECT id, user_id, name, color, period, total_amount FROM fn_category_amount_per_period($1) WHERE " + f.Where + f.OrderBy + f.Limit + f.Offset
	args := append([]any{period}, f.Args...)

	rows, err := db.QueryContext(ctx, sql, args...)

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
