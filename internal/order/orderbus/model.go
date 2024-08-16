package orderbus

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

// =============================================================================

// Product represents information about an individual product.
type Product struct {
	ID          uuid.UUID
	Name        Name
	Description string
	ImageURL    url.URL
	Price       int64
	Quantity    int32
	DateCreated time.Time
	DateUpdated time.Time
}

// =============================================================================

// NewProduct contains information needed to create a new product.
type NewProduct struct {
	Name        Name
	Description string
	ImageURL    url.URL
	Price       int64
	Quantity    int32
}

// =============================================================================

// UpdateProduct contains information needed to update a product.
type UpdateProduct struct {
	Name        *Name
	Description *string
	ImageURL    *url.URL
	Price       *int64
	Quantity    *int32
}
