package orderdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

func (s *Store) QueryByID(ctx context.Context, ordID uuid.UUID) (orderbus.Order, error) {
	// ==========================================================
	// Get order

	data := struct {
		ID string `db:"order_id"`
	}{
		ID: ordID.String(),
	}

	const q = `
	SELECT
		order_id, user_id, amount, status, date_created, date_updated
	FROM
		orders
	WHERE 
		order_id = :order_id`

	var dbOrd order
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbOrd); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return orderbus.Order{}, fmt.Errorf("db: %w", orderbus.ErrNotFound)
		}
		return orderbus.Order{}, fmt.Errorf("db: %w", err)
	}

	// ==========================================================
	// Get order items

	const itmQ = `
	SELECT
		order_item_id, order_id, product_id, price, quantity, date_created, date_updated
	FROM
		order_items
	WHERE
		order_id = :order_id`

	var dbItms []orderItem
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, itmQ, data, &dbItms); err != nil {
		return orderbus.Order{}, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusOrder(dbOrd, dbItms)
}
