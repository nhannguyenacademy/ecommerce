package productdb

import (
	"fmt"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
)

var sortByFields = map[string]string{
	productbus.SortByDateCreated: "date_created",
	productbus.SortByName:        "name",
	productbus.SortByPrice:       "price",
	productbus.SortByQuantity:    "quantity",
}

func orderByClause(sortBy sort.By) (string, error) {
	by, exists := sortByFields[sortBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", sortBy.Field)
	}

	return " ORDER BY " + by + " " + sortBy.Direction, nil
}
