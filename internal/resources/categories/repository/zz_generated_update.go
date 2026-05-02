// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *CategoriesRepoImpl) Update(ctx context.Context, db util.Executer, data UpdateCategoryDTO) error {
	filter := querybuilder.Get(ctx)
	var sets []string
	var values []interface{}
	if data.Name.Set {
		values = append(values, data.Name.Value)
		sets = append(sets, fmt.Sprintf("name = $%d", len(values)))
	}
	if data.Color.Set {
		values = append(values, data.Color.Value)
		sets = append(sets, fmt.Sprintf("color = $%d", len(values)))
	}

	f := filter.ToSQL(len(values) + 1)
	sql := "UPDATE categories SET " + strings.Join(sets, ", ") + " WHERE " + f.Where
	values = append(values, f.Args...)

	_, err := db.ExecContext(ctx, sql, values...)

	return err
}
