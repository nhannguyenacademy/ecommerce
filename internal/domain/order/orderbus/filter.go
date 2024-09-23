package orderbus

import (
	"github.com/google/uuid"
	"time"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
	UserID           *uuid.UUID
	Status           *Status
}
