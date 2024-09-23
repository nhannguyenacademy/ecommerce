package orderbus

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

// =============================================================================

type Order struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Amount      int64
	Status      Status
	DateCreated time.Time
	DateUpdated time.Time
}

type OrderWithItems struct {
	Order
	Items []OrderItem
}

// =============================================================================

type OrderItem struct {
	ID              uuid.UUID
	OrderID         uuid.UUID
	ProductID       uuid.UUID
	ProductName     string
	ProductImageURL url.URL
	Price           int64
	Quantity        int32
	DateCreated     time.Time
	DateUpdated     time.Time
}

// =============================================================================

type NewOrder struct {
	UserID uuid.UUID
	Items  []NewOrderItem
}

// =============================================================================

type NewOrderItem struct {
	ProductID       uuid.UUID
	ProductName     string
	ProductImageURL url.URL
	Price           int64
	Quantity        int32
}
