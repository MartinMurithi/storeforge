package postgres

import (
	"context"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"

)

type Adapter struct {
	pool *Pool
}

func New(p *Pool) database.DB {
    return &Adapter{pool: p}
}

func (a *Adapter) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) database.Row {
	return a.pool.QueryRow(ctx, sql, args...)
}

func (a *Adapter) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (database.CommandTag, error) {
	return a.pool.Exec(ctx, sql, args...)
}


func (a *Adapter) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
    return a.pool.Query(ctx, sql, args...)
}
