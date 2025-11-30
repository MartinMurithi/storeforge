// WrapDbError is a universal database error handler.
//
// # What it does
// 1. **Checks context cancellation FIRST** (timeout/client cancel)
// 2. **Maps Postgres errors** to domain errors (23505 → ErrUserAlreadyExists)
// 3. **Wraps generic errors** with operation context
//
// # `op` parameter
// Operation name for error messages and logging:
// - `"UserRepository.CreateUser"`
// - `"OrderRepository.CreateOrder"`
// - `"ReportRepository.GetMetrics"`
//
// # Usage (EVERY repository method)
//
// 	err := db.QueryRow(ctx, query...).Scan(&result)
// 	return WrapDbError(ctx, "UserRepository.CreateUser", 2*time.Second, err)
//
// # Examples
//
// | Scenario | Input | Output |
// |----------|-------|--------|
// | Success | `nil` | `nil` |
// | Timeout | `ctx.DeadlineExceeded` | `"UserRepository.CreateUser: timeout after 2s"` |
// | Postgres 23505 | `pgconn.PgError` | `dberrors.ErrUserAlreadyExists` |
// | Network | `"dial tcp"` | `"UserRepository.CreateUser: dial tcp"` |

package dbhelper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/pkg/dberrors"
	"github.com/jackc/pgx/v5/pgconn"
)

func WrapDbError(ctx context.Context, op string, timeout time.Duration, err error) error {

	// Check ctx error, even if it's nil
	ctxErr := ctx.Err()

	if ctxErr != nil {
		switch ctxErr {
		case context.DeadlineExceeded:
			return fmt.Errorf("%s: timeout after %v %w", op, timeout, ctxErr)
		case context.Canceled:
			return fmt.Errorf("%s operation cancelled by client %w", op, ctxErr)
		}
	}

	//check db errors
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, pgErr) {
			return dberrors.MapPostgresError(err)
		}
		return fmt.Errorf("%s %w", op, err)
	}

	return nil

}
