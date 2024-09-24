package orderbus

import (
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sort"
)

// DefaultSortBy represents the default way we sort.
var DefaultSortBy = sort.NewBy(SortByDateCreated, sort.DESC)

// Set of fields that the results can be ordered by.
const (
	SortByDateCreated = "date_created"
	SortByAmount      = "amount"
	SortByStatus      = "status"
)
