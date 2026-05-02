// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/util"
)

func (r *UsersRepoImpl) Insert(ctx context.Context, db util.Executer, data CreateUserDTO) error {
	var columns []string
	var values []interface{}
	var placeholders []string
	columns = append(columns, "id")
	values = append(values, data.ID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "name")
	values = append(values, data.Name)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "email")
	values = append(values, data.Email)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "avatar_url")
	values = append(values, data.AvatarURL)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	columns = append(columns, "username")
	values = append(values, data.Username)
	placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))

	sql := fmt.Sprintf("INSERT INTO users (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err := db.ExecContext(ctx, sql, values...)

	return err
}
