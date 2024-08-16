package orderbus

import "github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByDateCreated, order.DESC)

// Set of fields that the results can be ordered by.
const (
	OrderByDateCreated = "date_created"
	OrderByName        = "name"
	OrderByPrice       = "price"
	OrderByQuantity    = "quantity"
)
