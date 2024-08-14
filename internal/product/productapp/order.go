package productapp

import (
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
)

// Orders

var defaultOrderBy = order.NewBy("date_created", order.ASC)

var orderByFields = map[string]string{
	"name":         productbus.OrderByName,
	"date_created": productbus.OrderByDateCreated,
	"price":        productbus.OrderByPrice,
	"quantity":     productbus.OrderByQuantity,
}
