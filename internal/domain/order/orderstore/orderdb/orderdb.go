// Package orderdb provides the set of APIs for database access.
package orderdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sort"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

// Store manages the set of APIs for database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (orderbus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// ========================================================
// Orders

func (s *Store) Create(ctx context.Context, order orderbus.Order) error {
	const ordQ = `
	INSERT INTO orders
		(order_id, user_id, amount, status, date_created, date_updated)
	VALUES
		(:order_id, :user_id, :amount, :status, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, ordQ, toDBOrder(order)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, order orderbus.Order) error {
	const q = `
	DELETE FROM
		orders
	WHERE
		order_id = :order_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBOrder(order)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Query(ctx context.Context, filter orderbus.QueryFilter, sortBy sort.By, page page.Page) ([]orderbus.Order, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		order_id, user_id, amount, status, date_created, date_updated
	FROM
		orders`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(sortBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var rows []orderRow
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &rows); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusOrders(rows)
}

func (s *Store) Count(ctx context.Context, filter orderbus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		orders`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

func (s *Store) QueryByID(ctx context.Context, orderID uuid.UUID) (orderbus.Order, error) {
	// ==========================================================
	// Get order

	data := struct {
		ID string `db:"order_id"`
	}{
		ID: orderID.String(),
	}

	const q = `
	SELECT
		order_id, user_id, amount, status, date_created, date_updated
	FROM
		orders
	WHERE 
		order_id = :order_id`

	var row orderRow
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &row); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return orderbus.Order{}, fmt.Errorf("db: %w", orderbus.ErrNotFound)
		}
		return orderbus.Order{}, fmt.Errorf("db: %w", err)
	}

	return toBusOrder(row)
}

func (s *Store) UpdateStatus(ctx context.Context, order orderbus.Order, status orderbus.Status) error {
	data := struct {
		OrderID   uuid.UUID `db:"order_id"`
		Status    string    `db:"status"`
		NewStatus string    `db:"new_status"`
	}{
		OrderID:   order.ID,
		Status:    order.Status.String(),
		NewStatus: status.String(),
	}

	const q = `
	UPDATE orders
	SET status = :new_status
	WHERE order_id = :order_id AND status = :status`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// ========================================================
// Order Items

func (s *Store) QueryOrderItems(ctx context.Context, order orderbus.Order) ([]orderbus.OrderItem, error) {
	const itmQ = `
	SELECT
		order_item_id, order_id, product_id, product_name, product_image_url, price, quantity, date_created, date_updated
	FROM
		order_items
	WHERE
		order_id = :order_id`

	var rows []orderItemRow
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, itmQ, toDBOrder(order), &rows); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusOrderItems(rows)
}

func (s *Store) DeleteOrderItems(ctx context.Context, order orderbus.Order) error {
	const q = `
	DELETE FROM
		order_items
	WHERE
		order_id = :order_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBOrder(order)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) CreateOrderItems(ctx context.Context, items []orderbus.OrderItem) error {
	const ordItmQ = `
	INSERT INTO order_items
		(order_item_id, order_id, product_id, product_name, product_image_url, price, quantity, date_created, date_updated)
	VALUES
		(:order_item_id, :order_id, :product_id, :product_name, :product_image_url, :price, :quantity, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, ordItmQ, toDBOrderItems(items)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
