package orderdb

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/order/orderbus"
	"net/url"
	"time"
)

// ========================================================

type orderRow struct {
	ID          uuid.UUID `db:"order_id"`
	UserID      uuid.UUID `db:"user_id"`
	Amount      int64     `db:"amount"`
	Status      string    `db:"status"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBOrder(bus orderbus.Order) orderRow {
	return orderRow{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Amount:      bus.Amount,
		Status:      bus.Status.String(),
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}
}

func toBusOrder(row orderRow) (orderbus.Order, error) {
	orderStatus, err := orderbus.ParseStatus(row.Status)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse status: %w", err)
	}

	bus := orderbus.Order{
		ID:          row.ID,
		UserID:      row.UserID,
		Amount:      row.Amount,
		Status:      orderStatus,
		DateCreated: row.DateCreated.UTC(),
		DateUpdated: row.DateUpdated.UTC(),
	}

	return bus, nil
}

func toBusOrders(rows []orderRow) ([]orderbus.Order, error) {
	orders := make([]orderbus.Order, len(rows))
	for i, row := range rows {
		ord, err := toBusOrder(row)
		if err != nil {
			return nil, fmt.Errorf("to bus order: %w", err)
		}

		orders[i] = ord
	}

	return orders, nil
}

// ========================================================

type orderItemRow struct {
	ID              uuid.UUID `db:"order_item_id"`
	OrderID         uuid.UUID `db:"order_id"`
	ProductID       uuid.UUID `db:"product_id"`
	ProductName     string    `db:"product_name"`
	ProductImageURL string    `db:"product_image_url"`
	Price           int64     `db:"price"`
	Quantity        int32     `db:"quantity"`
	DateCreated     time.Time `db:"date_created"`
	DateUpdated     time.Time `db:"date_updated"`
}

func toDBOrderItem(bus orderbus.OrderItem) orderItemRow {
	return orderItemRow{
		ID:              bus.ID,
		OrderID:         bus.OrderID,
		ProductID:       bus.ProductID,
		ProductName:     bus.ProductName,
		ProductImageURL: bus.ProductImageURL.String(),
		Price:           bus.Price,
		Quantity:        bus.Quantity,
		DateCreated:     bus.DateCreated.UTC(),
		DateUpdated:     bus.DateUpdated.UTC(),
	}
}

func toDBOrderItems(items []orderbus.OrderItem) []orderItemRow {
	rows := make([]orderItemRow, len(items))
	for i, item := range items {
		rows[i] = toDBOrderItem(item)
	}

	return rows
}

func toBusOrderItem(row orderItemRow) (orderbus.OrderItem, error) {
	productImageURL, err := url.Parse(row.ProductImageURL)
	if err != nil {
		return orderbus.OrderItem{}, fmt.Errorf("parse product image url: %w", err)
	}

	item := orderbus.OrderItem{
		ID:              row.ID,
		ProductID:       row.ProductID,
		ProductName:     row.ProductName,
		ProductImageURL: *productImageURL,
		Price:           row.Price,
		Quantity:        row.Quantity,
		DateCreated:     row.DateCreated.UTC(),
		DateUpdated:     row.DateUpdated.UTC(),
	}

	return item, nil
}

func toBusOrderItems(rows []orderItemRow) ([]orderbus.OrderItem, error) {
	items := make([]orderbus.OrderItem, len(rows))
	for i, row := range rows {
		item, err := toBusOrderItem(row)
		if err != nil {
			return nil, fmt.Errorf("to bus order item: %w", err)
		}

		items[i] = item
	}

	return items, nil
}
