package orderdb

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"time"
)

type order struct {
	ID          uuid.UUID `db:"order_id"`
	UserID      uuid.UUID `db:"user_id"`
	Amount      int64     `db:"amount"`
	Status      string    `db:"status"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type orderItem struct {
	ID          uuid.UUID `db:"order_item_id"`
	OrderID     uuid.UUID `db:"order_id"`
	ProductID   uuid.UUID `db:"product_id"`
	Price       int64     `db:"price"`
	Quantity    int32     `db:"quantity"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBOrder(bus orderbus.Order) (order, []orderItem) {
	itms := make([]orderItem, len(bus.Items))
	for i, itm := range bus.Items {
		itms[i] = toDBOrderItem(itm)
	}

	ord := order{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Amount:      bus.Amount,
		Status:      bus.Status.String(),
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}

	return ord, itms
}

func toBusOrder(db order, itms []orderItem) (orderbus.Order, error) {
	ordStatus, err := orderbus.ParseStatus(db.Status)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse status: %w", err)
	}

	itmsBus := make([]orderbus.OrderItem, len(itms))
	for i, itm := range itms {
		itmBus, err := toBusOrderItem(itm)
		if err != nil {
			return orderbus.Order{}, fmt.Errorf("to bus order item: %w", err)
		}
		itmsBus[i] = itmBus
	}

	bus := orderbus.Order{
		ID:          db.ID,
		UserID:      db.UserID,
		Amount:      db.Amount,
		Status:      ordStatus,
		DateCreated: db.DateCreated.UTC(),
		DateUpdated: db.DateUpdated.UTC(),
		Items:       itmsBus,
	}

	return bus, nil
}

func toBusOrders(dbOrds []order, dbOrdItms []orderItem) ([]orderbus.Order, error) {
	itemsMap := make(map[uuid.UUID][]orderItem)
	for _, item := range dbOrdItms {
		itemsMap[item.OrderID] = append(itemsMap[item.OrderID], item)
	}

	ords := make([]orderbus.Order, len(dbOrds))
	for i, dbOrd := range dbOrds {
		itms := itemsMap[dbOrd.ID]

		ord, err := toBusOrder(dbOrd, itms)
		if err != nil {
			return nil, fmt.Errorf("to bus order: %w", err)
		}

		ords[i] = ord
	}

	return ords, nil
}

func toDBOrderItem(bus orderbus.OrderItem) orderItem {
	return orderItem{
		ID:          bus.ID,
		OrderID:     bus.OrderID,
		ProductID:   bus.ProductID,
		Price:       bus.Price,
		Quantity:    bus.Quantity,
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}
}

func toBusOrderItem(db orderItem) (orderbus.OrderItem, error) {
	itm := orderbus.OrderItem{
		ID:          db.ID,
		ProductID:   db.ProductID,
		Price:       db.Price,
		Quantity:    db.Quantity,
		DateCreated: db.DateCreated.UTC(),
		DateUpdated: db.DateUpdated.UTC(),
	}

	return itm, nil
}
