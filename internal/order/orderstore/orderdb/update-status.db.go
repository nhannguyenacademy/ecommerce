package orderdb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

func (s *Store) UpdateStatus(ctx context.Context, ord orderbus.Order, status orderbus.Status) error {
	uo := struct {
		OrderID   uuid.UUID `db:"order_id"`
		Status    string    `db:"status"`
		NewStatus string    `db:"new_status"`
	}{
		OrderID:   ord.ID,
		Status:    ord.Status.String(),
		NewStatus: status.String(),
	}

	const q = `
	UPDATE orders
	SET status = :new_status
	WHERE order_id = :order_id AND status = :status`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, uo); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
