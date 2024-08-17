package orderdb

import (
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

func (s *Store) Create(ctx context.Context, ord orderbus.Order) error {
	o, itms := toDBOrder(ord)

	const ordQ = `
	INSERT INTO orders
		(order_id, user_id, amount, status, date_created, date_updated)
	VALUES
		(:order_id, :user_id, :amount, :status, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, ordQ, o); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	const ordItmQ = `
	INSERT INTO order_items
		(order_item_id, order_id, product_id, price, quantity, date_created, date_updated)
	VALUES
		(:order_item_id, :order_id, :product_id, :price, :quantity, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, ordItmQ, itms); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
