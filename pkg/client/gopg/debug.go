package gopg

import (
	"context"

	"github.com/boxgo/box/pkg/logger"
	"github.com/go-pg/pg/v10"
)

// DebugHook is a query hook that logs an error with a query if there are any.
// It can be installed with:
//
//   db.AddQueryHook(pgext.DebugHook{})
// Copy from https://github.com/go-pg/pg/tree/v10/extra/pgdebug
type DebugHook struct {
	// Verbose causes hook to print all queries (even those without an error).
	Verbose bool
}

var _ pg.QueryHook = (*DebugHook)(nil)

func (h DebugHook) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if evt.Err != nil {
		logger.Trace(ctx).Errorw("PgExecuteError", "error", evt.Err)
	} else if h.Verbose {
		logger.Trace(ctx).Infow("PgExecuteStart", "sql", string(q))
	}

	return ctx, nil
}

func (h DebugHook) AfterQuery(ctx context.Context, evt *pg.QueryEvent) error {
	if evt.Err != nil {
		logger.Trace(ctx).Errorw("PGExecuteError", "error", evt.Err)
	} else if h.Verbose {
		logger.Trace(ctx).Infow("PgExecuteEnd", "affected", evt.Result.RowsAffected(), "returned", evt.Result.RowsReturned())
	}

	return nil
}
