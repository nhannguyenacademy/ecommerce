package productbus

import (
	"github.com/google/uuid"
	"net/mail"
	"time"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID               *uuid.UUID
	Name             *Name
	Email            *mail.Address
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}
