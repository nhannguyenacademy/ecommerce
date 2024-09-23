package orderdb

import (
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
)

var sortByFields = map[string]string{
	orderbus.SortByDateCreated: "date_created",
	orderbus.SortByAmount:      "amount",
	orderbus.SortByStatus:      "status",
}

func orderByClause(sortBy sort.By) (string, error) {
	by, exists := sortByFields[sortBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", sortBy.Field)
	}

	return " ORDER BY " + by + " " + sortBy.Direction, nil
}
