package repository

import (
	"context"
	"database/sql"
)

// sqlExecutor abstracts *sql.DB / *sql.Tx for repositories that run raw SQL.
type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
