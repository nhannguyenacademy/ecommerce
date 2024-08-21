package orderdb

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"net/url"
	"time"
)

// ========================================================

type order struct {
	ID          uuid.UUID `db:"order_id"`
	UserID      uuid.UUID `db:"user_id"`
	Amount      int64     `db:"amount"`
	Status      string    `db:"status"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBOrder(bus orderbus.Order) order {
	ord := order{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Amount:      bus.Amount,
		Status:      bus.Status.String(),
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}

	return ord
}

func toBusOrder(db order) (orderbus.Order, error) {
	ordStatus, err := orderbus.ParseStatus(db.Status)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse status: %w", err)
	}

	bus := orderbus.Order{
		ID:          db.ID,
		UserID:      db.UserID,
		Amount:      db.Amount,
		Status:      ordStatus,
		DateCreated: db.DateCreated.UTC(),
		DateUpdated: db.DateUpdated.UTC(),
	}

	return bus, nil
}

func toBusOrders(dbOrds []order) ([]orderbus.Order, error) {
	ords := make([]orderbus.Order, len(dbOrds))
	for i, dbOrd := range dbOrds {
		ord, err := toBusOrder(dbOrd)
		if err != nil {
			return nil, fmt.Errorf("to bus order: %w", err)
		}

		ords[i] = ord
	}

	return ords, nil
}

// ========================================================

type orderItem struct {
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

func toDBOrderItem(bus orderbus.OrderItem) orderItem {
	return orderItem{
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

func toDBOrderItems(busItms []orderbus.OrderItem) []orderItem {
	itms := make([]orderItem, len(busItms))
	for i, busItm := range busItms {
		itms[i] = toDBOrderItem(busItm)
	}

	return itms
}

func toBusOrderItem(db orderItem) (orderbus.OrderItem, error) {
	productImageURL, err := url.Parse(db.ProductImageURL)
	if err != nil {
		return orderbus.OrderItem{}, fmt.Errorf("parse product image url: %w", err)
	}

	itm := orderbus.OrderItem{
		ID:              db.ID,
		ProductID:       db.ProductID,
		ProductName:     db.ProductName,
		ProductImageURL: *productImageURL,
		Price:           db.Price,
		Quantity:        db.Quantity,
		DateCreated:     db.DateCreated.UTC(),
		DateUpdated:     db.DateUpdated.UTC(),
	}

	return itm, nil
}

func toBusOrderItems(dbItems []orderItem) ([]orderbus.OrderItem, error) {
	items := make([]orderbus.OrderItem, len(dbItems))
	for i, dbItem := range dbItems {
		item, err := toBusOrderItem(dbItem)
		if err != nil {
			return nil, fmt.Errorf("to bus order item: %w", err)
		}

		items[i] = item
	}

	return items, nil
}
