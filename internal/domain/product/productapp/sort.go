package productapp

import (
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
)

var defaultSortBy = sort.NewBy("date_created", sort.ASC)

var sortByFields = map[string]string{
	"name":         productbus.SortByName,
	"date_created": productbus.SortByDateCreated,
	"price":        productbus.SortByPrice,
	"quantity":     productbus.SortByQuantity,
}
