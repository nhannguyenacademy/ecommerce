package orderdb

import (
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	ordering "github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
)

var orderByFields = map[string]string{
	orderbus.OrderByDateCreated: "date_created",
	orderbus.OrderByAmount:      "amount",
	orderbus.OrderByStatus:      "status",
}

func orderByClause(orderBy ordering.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
