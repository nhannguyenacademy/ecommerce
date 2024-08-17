package orderdb

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	ordering "github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
)

func (s *Store) Query(
	ctx context.Context,
	filter orderbus.QueryFilter,
	orderBy ordering.By,
	page page.Page,
) ([]orderbus.Order, error) {
	// ==========================================================
	// Get orders

	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const ordQ = `
	SELECT
		order_id, user_id, amount, status, date_created, date_updated
	FROM
		orders`

	buf := bytes.NewBufferString(ordQ)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbOrds []order
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbOrds); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	// ==========================================================
	// Get order items

	const itmQ = `
	SELECT
		order_item_id, order_id, product_id, price, quantity, date_created, date_updated
	FROM
		order_items
	WHERE
		order_id IN (:order_ids)`

	ordIDs := make([]string, len(dbOrds))
	for i, dbOrd := range dbOrds {
		ordIDs[i] = dbOrd.ID.String()
	}

	inData := map[string]any{
		"order_ids": ordIDs,
	}
	var dbItms []orderItem
	if err := sqldb.NamedQuerySliceUsingIn(ctx, s.log, s.db, itmQ, inData, &dbItms); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusOrders(dbOrds, dbItms)
}
