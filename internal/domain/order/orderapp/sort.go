package orderapp

import (
	"github.com/nhannguyenacademy/ecommerce/internal/domain/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sort"
)

var defaultSortBy = sort.NewBy("date_created", sort.ASC)

var sortByFields = map[string]string{
	"date_created": orderbus.SortByDateCreated,
	"amount":       orderbus.SortByAmount,
	"status":       orderbus.SortByStatus,
}
