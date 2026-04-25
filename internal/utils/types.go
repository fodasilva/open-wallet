package utils

import (
	"context"
	"database/sql"

	"github.com/felipe1496/open-wallet/internal/utils/querybuilder"
)

type Executer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type ResponseData[T any] struct {
	Data T `json:"data"`
}

type PaginatedResponse[T any] struct {
	Data  T                     `json:"data"`
	Query querybuilder.Metadata `json:"query"`
}

type OptionalNullable[T any] struct {
	Set   bool
	Value *T
}
