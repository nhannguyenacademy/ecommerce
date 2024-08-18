package orderapp

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"net/http"
	"time"
)

// =============================================================================

// queryParams represents the set of possible query strings.
type queryParams struct {
	Page             string
	Rows             string
	OrderBy          string
	StartCreatedDate string
	EndCreatedDate   string
	UserID           string
	Status           string
}

func parseQueryParams(r *http.Request) queryParams {
	values := r.URL.Query()

	filter := queryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("row"),
		OrderBy:          values.Get("order_by"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
		UserID:           values.Get("user_id"),
		Status:           values.Get("status"),
	}

	return filter
}

// ===================================================

type order struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	Amount      int64       `json:"amount"`
	Status      string      `json:"status"`
	DateCreated string      `json:"date_created"`
	DateUpdated string      `json:"date_updated"`
	Items       []orderItem `json:"items"`
}

func toAppOrder(bus orderbus.Order) order {
	return order{
		ID:          bus.ID.String(),
		UserID:      bus.UserID.String(),
		Amount:      bus.Amount,
		Status:      bus.Status.String(),
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
		Items:       toAppOrderItems(bus.Items),
	}
}

// ===================================================

type orderItem struct {
	ID          string `json:"id"`
	OrderID     string `json:"order_id"`
	ProductID   string `json:"product_id"`
	Price       int64  `json:"price"`
	Quantity    int32  `json:"quantity"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

func toAppOrderItem(bus orderbus.OrderItem) orderItem {
	return orderItem{
		ID:          bus.ID.String(),
		OrderID:     bus.OrderID.String(),
		ProductID:   bus.ProductID.String(),
		Price:       bus.Price,
		Quantity:    bus.Quantity,
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppOrderItems(bus []orderbus.OrderItem) []orderItem {
	items := make([]orderItem, len(bus))
	for i, v := range bus {
		items[i] = toAppOrderItem(v)
	}
	return items
}

// ===================================================

type newOrder struct {
	UserID string         `json:"user_id"`
	Items  []newOrderItem `json:"items"`
}

type newOrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

func toBusNewOrder(app newOrder) (orderbus.NewOrder, error) {
	userID, err := uuid.Parse(app.UserID)
	if err != nil {
		return orderbus.NewOrder{}, fmt.Errorf("parsing user id: %w", err)
	}

	items, err := toBusNewOrderItems(app.Items)
	if err != nil {
		return orderbus.NewOrder{}, err
	}

	return orderbus.NewOrder{
		UserID: userID,
		Items:  items,
	}, nil
}

func toBusNewOrderItems(app []newOrderItem) ([]orderbus.NewOrderItem, error) {
	items := make([]orderbus.NewOrderItem, len(app))
	var err error
	for i, v := range app {
		items[i], err = toBusNewOrderItem(v)
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func toBusNewOrderItem(app newOrderItem) (orderbus.NewOrderItem, error) {
	productID, err := uuid.Parse(app.ProductID)
	if err != nil {
		return orderbus.NewOrderItem{}, fmt.Errorf("parsing product id: %w", err)
	}

	return orderbus.NewOrderItem{
		ProductID: productID,
		Quantity:  app.Quantity,
	}, nil
}
