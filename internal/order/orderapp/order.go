package orderapp

import (
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	ordering "github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
)

var defaultOrderBy = ordering.NewBy("date_created", ordering.ASC)

var orderByFields = map[string]string{
	"date_created": orderbus.OrderByDateCreated,
	"amount":       orderbus.OrderByAmount,
	"status":       orderbus.OrderByStatus,
}
