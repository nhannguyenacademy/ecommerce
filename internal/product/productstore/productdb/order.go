package productdb

import (
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
)

var orderByFields = map[string]string{
	productbus.OrderByDateCreated: "date_created",
	productbus.OrderByName:        "name",
	productbus.OrderByPrice:       "price",
	productbus.OrderByQuantity:    "quantity",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
