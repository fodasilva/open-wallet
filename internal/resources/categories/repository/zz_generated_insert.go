// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/util"
)

func (r *CategoriesRepoImpl) Insert(ctx context.Context, db util.Executer, data CreateCategoryDTO) error {
	var columns []string
	var values []interface{}
	var placeholders []string
	columns = append(columns, "id")
	values = append(values, data.ID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "user_id")
	values = append(values, data.UserID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "name")
	values = append(values, data.Name)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "color")
	values = append(values, data.Color)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))

	sql := fmt.Sprintf("INSERT INTO categories (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err := db.ExecContext(ctx, sql, values...)

	return err
}
