// Code generated. DO NOT EDIT.

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/felipe1496/open-wallet/internal/util"
	"github.com/felipe1496/open-wallet/internal/util/querybuilder"
)

func (r *UsersRepoImpl) Update(ctx context.Context, db util.Executer, data UpdateUserDTO) error {
	filter := querybuilder.Get(ctx)
	var sets []string
	var values []interface{}
	if data.Name.Set {
		values = append(values, data.Name.Value)
		sets = append(sets, fmt.Sprintf("name = $%d", len(values)))
	}
	if data.Email.Set {
		values = append(values, data.Email.Value)
		sets = append(sets, fmt.Sprintf("email = $%d", len(values)))
	}
	if data.AvatarURL.Set {
		values = append(values, data.AvatarURL.Value)
		sets = append(sets, fmt.Sprintf("avatar_url = $%d", len(values)))
	}
	if data.Username.Set {
		values = append(values, data.Username.Value)
		sets = append(sets, fmt.Sprintf("username = $%d", len(values)))
	}

	f := filter.ToSQL(len(values) + 1)
	sql := "UPDATE users SET " + strings.Join(sets, ", ") + " WHERE " + f.Where
	values = append(values, f.Args...)

	_, err := db.ExecContext(ctx, sql, values...)

	return err
}
