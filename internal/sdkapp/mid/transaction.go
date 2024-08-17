package mid

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

// BeginCommitRollback executes the transaction middleware functionality.
func BeginCommitRollback(l *logger.Logger, bgn sqldb.Beginner) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		hasCommitted := false

		l.Info(ctx, "BEGIN TRANSACTION")
		tx, err := bgn.Begin()
		if err != nil {
			response.Send(c, l, nil, errs.Newf(errs.Internal, "BEGIN TRANSACTION: %s", err))
			return
		}

		defer func() {
			if !hasCommitted {
				l.Info(ctx, "ROLLBACK TRANSACTION")
			}

			if err := tx.Rollback(); err != nil {
				if errors.Is(err, sql.ErrTxDone) {
					return
				}
				l.Info(ctx, "ROLLBACK TRANSACTION", "ERROR", err)
			}
		}()

		ctx = setTran(ctx, tx)

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if c.IsAborted() || len(c.Errors) > 0 {
			return
		}

		l.Info(ctx, "COMMIT TRANSACTION")
		if err := tx.Commit(); err != nil {
			response.Send(c, l, nil, errs.Newf(errs.Internal, "COMMIT TRANSACTION: %s", err))
			return
		}

		hasCommitted = true
	}
}
