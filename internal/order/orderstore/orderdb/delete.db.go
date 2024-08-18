package orderdb

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

func (s *Store) Delete(ctx context.Context, ord orderbus.Order) error {
	// this method needs to be run within a transaction
	if _, ok := s.db.(sqldb.CommitRollbacker); !ok {
		return fmt.Errorf("store: db not a CommitRollbacker")
	}

	dbOrd, _ := toDBOrder(ord)

	// ===============================================================
	// Delete order items

	const q1 = `
	DELETE FROM
		order_items
	WHERE
		order_id = :order_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q1, dbOrd); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	// ===============================================================
	// Delete order

	const q = `
	DELETE FROM
		orders
	WHERE
		order_id = :order_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, dbOrd); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
