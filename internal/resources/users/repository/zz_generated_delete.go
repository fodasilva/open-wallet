// Code generated. DO NOT EDIT.

package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

func (r *UsersRepoImpl) Delete(db utils.Executer, filter *utils.QueryOptsBuilder) error {
	query := squirrel.Delete("users").
		PlaceholderFormat(squirrel.Dollar)

	query = utils.DeleteOptsToSquirrel(query, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = db.Exec(sql, args...)

	return err
}
