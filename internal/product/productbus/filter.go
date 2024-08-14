package productbus

import (
	"time"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	Name             *Name
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
	StartPrice       *int64
	EndPrice         *int64
}
